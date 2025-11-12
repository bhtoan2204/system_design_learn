package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	port := getEnv("API_PORT", "8081")

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/api/echo", echoHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: loggingMiddleware(mux),
	}

	log.Printf("socket_api listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("socket_api failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil && err != io.EOF {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	if body == nil {
		body = map[string]any{}
	}

	query := map[string]string{}
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "echo from socket_api",
		"method":  r.Method,
		"path":    r.URL.Path,
		"body":    body,
		"query":   query,
		"headers": map[string]string{
			"content-type": r.Header.Get("Content-Type"),
			"x-trace-id":   r.Header.Get("X-Trace-Id"),
		},
	})
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
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
