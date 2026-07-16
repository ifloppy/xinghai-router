package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type keyContext struct{ userID, keyID string }
type contextKey struct{}

func (s *Service) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Handle("POST /admin/users", s.admin(s.createUser))
	mux.Handle("GET /admin/users", s.admin(s.listUsers))
	mux.Handle("POST /admin/keys", s.admin(s.createKey))
	mux.Handle("GET /admin/keys", s.admin(s.listKeys))
	mux.Handle("POST /admin/keys/{id}/revoke", s.admin(s.revokeKey))
	mux.Handle("POST /admin/channels", s.admin(s.createChannel))
	mux.Handle("GET /admin/channels", s.admin(s.listChannels))
	mux.Handle("POST /admin/channels/{id}/status", s.admin(s.setChannelStatus))
	mux.Handle("GET /admin/request-logs", s.admin(s.listLogs))
	mux.Handle("GET /admin/pricing", s.admin(s.listPricing))
	mux.Handle("POST /admin/pricing", s.admin(s.upsertPricing))
	mux.Handle("GET /admin/audit-logs", s.admin(s.listAuditLogs))
	mux.Handle("POST /admin/wallets/adjustments", s.admin(s.adjustBalance))
	mux.Handle("GET /admin/model-routes", s.admin(s.listModelRoutes))
	mux.Handle("POST /admin/model-routes", s.admin(s.createModelRoute))
	mux.Handle("POST /admin/quota-limits", s.admin(s.upsertQuota))
	mux.Handle("GET /me", s.api(s.me))
	mux.Handle("GET /me/keys", s.api(s.myKeys))
	mux.Handle("GET /me/usage", s.api(s.myUsage))
	mux.Handle("GET /me/ledger", s.api(s.myLedger))
	mux.Handle("GET /v1/models", s.api(s.models))
	mux.Handle("POST /v1/chat/completions", s.api(s.chatCompletions))
	return s.requestID(mux)
}
func (s *Service) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := randomID()
		if err != nil {
			writeError(w, 500, "internal_error", "could not create request id")
			return
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id)))
	})
}

type requestIDKey struct{}

func requestID(ctx context.Context) string { id, _ := ctx.Value(requestIDKey{}).(string); return id }
func (s *Service) admin(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !equalSecret(bearer(r), s.cfg.AdminToken) {
			writeError(w, 401, "unauthorized", "administrator token required")
			return
		}
		next(w, r)
	})
}
func (s *Service) api(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearer(r)
		if token == "" {
			writeError(w, 401, "invalid_api_key", "API key required")
			return
		}
		var k keyContext
		err := s.db.QueryRow(r.Context(), `select k.user_id,k.id from api_keys k join users u on u.id=k.user_id where k.secret_hash=$1 and k.revoked_at is null and (k.expires_at is null or k.expires_at>now()) and u.enabled`, hashSecret(token)).Scan(&k.userID, &k.keyID)
		if err != nil {
			writeError(w, 401, "invalid_api_key", "invalid or expired API key")
			return
		}
		if !s.limiter.allow(k.keyID) {
			writeError(w, 429, "rate_limit_exceeded", "too many requests")
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, k)))
	})
}

func (s *Service) me(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(contextKey{}).(keyContext)
	var email, name, role string
	var balance, reserved any
	err := s.db.QueryRow(r.Context(), `select u.email,u.name,u.role,coalesce(w.balance,0),coalesce(w.reserved,0) from users u left join user_wallets w on w.user_id=u.id where u.id=$1`, key.userID).Scan(&email, &name, &role, &balance, &reserved)
	if err != nil {
		writeError(w, 500, "internal_error", "could not load account")
		return
	}
	writeJSON(w, 200, map[string]any{"user_id": key.userID, "key_id": key.keyID, "email": email, "name": name, "role": role, "balance": balance, "reserved": reserved})
}
func (s *Service) myKeys(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(contextKey{}).(keyContext)
	rows, err := s.db.Query(r.Context(), `select id,name,key_prefix,expires_at,revoked_at,last_used_at,created_at from api_keys where user_id=$1 order by created_at desc`, key.userID)
	if err != nil {
		writeError(w, 500, "internal_error", "query failed")
		return
	}
	defer rows.Close()
	data := []map[string]any{}
	for rows.Next() {
		var id, name, prefix string
		var expires, revoked, used, created any
		if rows.Scan(&id, &name, &prefix, &expires, &revoked, &used, &created) == nil {
			data = append(data, map[string]any{"id": id, "name": name, "key_prefix": prefix, "expires_at": expires, "revoked_at": revoked, "last_used_at": used, "created_at": created})
		}
	}
	writeJSON(w, 200, map[string]any{"data": data})
}
func (s *Service) myUsage(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(contextKey{}).(keyContext)
	rows, err := s.db.Query(r.Context(), `select request_id,model,prompt_tokens,cached_prompt_tokens,completion_tokens,cost,status,created_at from usage_records where user_id=$1 order by created_at desc limit 100`, key.userID)
	if err != nil {
		writeError(w, 500, "internal_error", "query failed")
		return
	}
	defer rows.Close()
	data := []map[string]any{}
	for rows.Next() {
		var requestID, model, status string
		var prompt, cached, completion int
		var cost, created any
		if rows.Scan(&requestID, &model, &prompt, &cached, &completion, &cost, &status, &created) == nil {
			data = append(data, map[string]any{"request_id": requestID, "model": model, "prompt_tokens": prompt, "cached_prompt_tokens": cached, "completion_tokens": completion, "cost": cost, "status": status, "created_at": created})
		}
	}
	writeJSON(w, 200, map[string]any{"data": data})
}
func (s *Service) myLedger(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(contextKey{}).(keyContext)
	rows, err := s.db.Query(r.Context(), `select id,amount,balance_after,kind,request_id,note,created_at from wallet_ledger where user_id=$1 order by created_at desc limit 100`, key.userID)
	if err != nil {
		writeError(w, 500, "internal_error", "query failed")
		return
	}
	defer rows.Close()
	data := []map[string]any{}
	for rows.Next() {
		var id, kind, requestID, note string
		var amount, after, created any
		if rows.Scan(&id, &amount, &after, &kind, &requestID, &note, &created) == nil {
			data = append(data, map[string]any{"id": id, "amount": amount, "balance_after": after, "kind": kind, "request_id": requestID, "note": note, "created_at": created})
		}
	}
	writeJSON(w, 200, map[string]any{"data": data})
}
func bearer(r *http.Request) string {
	const p = "Bearer "
	v := r.Header.Get("Authorization")
	if strings.HasPrefix(v, p) {
		return strings.TrimSpace(strings.TrimPrefix(v, p))
	}
	return ""
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"message": msg, "type": code, "code": code}})
}
func decode(r *http.Request, target any) error {
	d := json.NewDecoder(io.LimitReader(r.Body, 2<<20))
	d.DisallowUnknownFields()
	return d.Decode(target)
}

var errInvalid = errors.New("invalid request")

func parseExpiry(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, value)
	return &t, err
}
