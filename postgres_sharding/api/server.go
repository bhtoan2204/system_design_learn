package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"database_sharding/persistent"

	"gorm.io/gorm"
)

type Server struct {
	db  *gorm.DB
	mux *http.ServeMux
}

func NewServer(db *gorm.DB) *Server {
	s := &Server{db: db, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.handleHealthz)
	s.mux.HandleFunc("/tenants", s.handleListTenants)
	s.mux.HandleFunc("/tenants/", s.handleGetTenant)
	s.mux.HandleFunc("/tenants/seed", s.handleSeedTenants)
	s.mux.HandleFunc("/tenants/rebalance", s.handleRebalanceTenants)
	s.mux.HandleFunc("/tenants/shards", s.handleTenantShardPlacements)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListTenants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}

	ctx := r.Context()
	limit, err := parsePositiveInt(r, "limit", 50, 1, 1000)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	offset, err := parsePositiveInt(r, "offset", 0, 0, 1_000_000_000)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	tenants, err := persistent.ListTenants(ctx, s.db, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	total, err := persistent.CountTenants(ctx, s.db)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"limit":   limit,
		"offset":  offset,
		"total":   total,
		"tenants": tenants,
	})
}

func (s *Server) handleSeedTenants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	if err := persistent.SeedTenants(r.Context(), s.db); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "seeding started/completed"})
}

func (s *Server) handleRebalanceTenants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	if err := persistent.RebalanceTenants(r.Context(), s.db); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "rebalance triggered"})
}

func (s *Server) handleTenantShardPlacements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	placements, err := persistent.ListTenantShardPlacements(r.Context(), s.db)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"count":      len(placements),
		"placements": placements,
	})
}

func (s *Server) handleGetTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	idStr := r.URL.Path[len("/tenants/"):]
	if idStr == "" {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid tenant id: %w", err))
		return
	}

	tenant, err := persistent.GetTenantWithShard(r.Context(), s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, tenant)
}

func parsePositiveInt(r *http.Request, key string, def, min, max int) (int, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return def, nil
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	if parsed < min {
		return 0, fmt.Errorf("%s must be >= %d", key, min)
	}
	if parsed > max {
		return 0, fmt.Errorf("%s must be <= %d", key, max)
	}
	return parsed, nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func methodNotAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
