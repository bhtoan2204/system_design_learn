package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type apiRequest struct {
	RequestID string            `json:"request_id"`
	Method    string            `json:"method"`
	Path      string            `json:"path"`
	Query     map[string]string `json:"query,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	Body      json.RawMessage   `json:"body,omitempty"`
}

type apiResponse struct {
	RequestID  string            `json:"request_id"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       json.RawMessage   `json:"body,omitempty"`
	Error      string            `json:"error,omitempty"`
}

type gateway struct {
	apiBaseURL string
	client     *http.Client
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	port := getEnv("GATEWAY_PORT", "8080")
	apiBaseURL := strings.TrimSuffix(getEnv("API_BASE_URL", "http://socket-api:8081"), "/")

	gw := &gateway{
		apiBaseURL: apiBaseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/ws", gw.handleWebSocket)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: loggingMiddleware(mux),
	}

	log.Printf("socket_gateway listening on :%s (api: %s)", port, apiBaseURL)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("socket_gateway failed: %v", err)
	}
}

func (g *gateway) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	for {
		var req apiRequest
		if err := conn.ReadJSON(&req); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket read error: %v", err)
			}
			break
		}

		resp := g.forwardRequest(r.Context(), req)
		if err := conn.WriteJSON(resp); err != nil {
			log.Printf("websocket write error: %v", err)
			break
		}
	}
}

func (g *gateway) forwardRequest(ctx context.Context, req apiRequest) apiResponse {
	resp := apiResponse{
		RequestID: req.RequestID,
	}

	if req.Path == "" {
		resp.StatusCode = http.StatusBadRequest
		resp.Error = "path is required"
		return resp
	}

	method := strings.ToUpper(req.Method)
	if method == "" {
		method = http.MethodGet
	}

	targetURL, err := g.buildURL(req.Path, req.Query)
	if err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Error = err.Error()
		return resp
	}

	var bodyReader io.Reader
	if len(req.Body) > 0 {
		bodyReader = bytes.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, targetURL, bodyReader)
	if err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Error = err.Error()
		return resp
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	if _, ok := req.Headers["Content-Type"]; !ok && len(req.Body) > 0 {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	httpResp, err := g.client.Do(httpReq)
	if err != nil {
		resp.StatusCode = http.StatusBadGateway
		resp.Error = err.Error()
		return resp
	}
	defer httpResp.Body.Close()

	resp.StatusCode = httpResp.StatusCode
	resp.Headers = map[string]string{
		"content-type": httpResp.Header.Get("Content-Type"),
	}

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	resp.Body = json.RawMessage(bodyBytes)
	return resp
}

func (g *gateway) buildURL(path string, query map[string]string) (string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	target := g.apiBaseURL + path
	parsed, err := url.Parse(target)
	if err != nil {
		return "", err
	}

	q := parsed.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	parsed.RawQuery = q.Encode()
	return parsed.String(), nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
