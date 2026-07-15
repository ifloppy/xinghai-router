package app

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Role  string `json:"role"`
	}
	if decode(r, &in) != nil || in.Email == "" || in.Name == "" {
		writeError(w, 400, "invalid_request", "email and name are required")
		return
	}
	if in.Role == "" {
		in.Role = "user"
	}
	if in.Role != "user" && in.Role != "operator" && in.Role != "admin" {
		writeError(w, 400, "invalid_request", "invalid role")
		return
	}
	id, _ := randomID()
	_, err := s.db.Exec(r.Context(), `insert into users(id,email,name,role) values($1,$2,$3,$4)`, id, strings.ToLower(in.Email), in.Name, in.Role)
	if err != nil {
		writeError(w, 409, "conflict", "email already exists")
		return
	}
	writeJSON(w, 201, map[string]any{"id": id, "email": strings.ToLower(in.Email), "name": in.Name, "role": in.Role})
}
func (s *Service) listUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(r.Context(), `select id,email,name,role,enabled,created_at from users order by created_at desc`)
	if err != nil {
		writeError(w, 500, "internal_error", "query failed")
		return
	}
	defer rows.Close()
	out := []map[string]any{}
	for rows.Next() {
		var id, email, name, role string
		var enabled bool
		var created any
		rows.Scan(&id, &email, &name, &role, &enabled, &created)
		out = append(out, map[string]any{"id": id, "email": email, "name": name, "role": role, "enabled": enabled, "created_at": created})
	}
	writeJSON(w, 200, map[string]any{"data": out})
}
func (s *Service) createKey(w http.ResponseWriter, r *http.Request) {
	var in struct {
		UserID    string `json:"user_id"`
		Name      string `json:"name"`
		ExpiresAt string `json:"expires_at"`
	}
	if decode(r, &in) != nil || in.UserID == "" || in.Name == "" {
		writeError(w, 400, "invalid_request", "user_id and name are required")
		return
	}
	expires, err := parseExpiry(in.ExpiresAt)
	if err != nil {
		writeError(w, 400, "invalid_request", "expires_at must be RFC3339")
		return
	}
	secret, err := randomSecret("sk-xh-")
	if err != nil {
		writeError(w, 500, "internal_error", "key generation failed")
		return
	}
	id, _ := randomID()
	_, err = s.db.Exec(r.Context(), `insert into api_keys(id,user_id,name,key_prefix,secret_hash,expires_at) values($1,$2,$3,$4,$5,$6)`, id, in.UserID, in.Name, secret[:12], hashSecret(secret), expires)
	if err != nil {
		writeError(w, 400, "invalid_request", "unknown user")
		return
	}
	writeJSON(w, 201, map[string]any{"id": id, "name": in.Name, "key": secret, "expires_at": expires})
}
func (s *Service) listKeys(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(r.Context(), `select id,user_id,name,key_prefix,expires_at,revoked_at,last_used_at,created_at from api_keys order by created_at desc`)
	if err != nil {
		writeError(w, 500, "internal_error", "query failed")
		return
	}
	defer rows.Close()
	data := []map[string]any{}
	for rows.Next() {
		var id, uid, name, prefix string
		var expiry, revoked, used, created any
		rows.Scan(&id, &uid, &name, &prefix, &expiry, &revoked, &used, &created)
		data = append(data, map[string]any{"id": id, "user_id": uid, "name": name, "key_prefix": prefix, "expires_at": expiry, "revoked_at": revoked, "last_used_at": used, "created_at": created})
	}
	writeJSON(w, 200, map[string]any{"data": data})
}
func (s *Service) revokeKey(w http.ResponseWriter, r *http.Request) {
	result, err := s.db.Exec(r.Context(), `update api_keys set revoked_at=coalesce(revoked_at, now()) where id=$1`, r.PathValue("id"))
	if err != nil || result.RowsAffected() != 1 {
		writeError(w, 404, "not_found", "API key not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (s *Service) createChannel(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name     string   `json:"name"`
		BaseURL  string   `json:"base_url"`
		APIKey   string   `json:"api_key"`
		Models   []string `json:"models"`
		Priority int      `json:"priority"`
	}
	if decode(r, &in) != nil || in.Name == "" || in.APIKey == "" || len(in.Models) == 0 {
		writeError(w, 400, "invalid_request", "name, api_key, and models are required")
		return
	}
	u, err := url.Parse(in.BaseURL)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		writeError(w, 400, "invalid_request", "base_url must be an HTTPS URL")
		return
	}
	encrypted, err := crypt(s.cfg.EncryptionKey, in.APIKey, false)
	if err != nil {
		writeError(w, 500, "internal_error", "credential encryption failed")
		return
	}
	models, _ := json.Marshal(in.Models)
	id, _ := randomID()
	_, err = s.db.Exec(r.Context(), `insert into channels(id,name,base_url,api_key,models,priority) values($1,$2,$3,$4,$5,$6)`, id, in.Name, strings.TrimRight(in.BaseURL, "/"), encrypted, models, in.Priority)
	if err != nil {
		writeError(w, 409, "conflict", "channel name already exists")
		return
	}
	writeJSON(w, 201, map[string]any{"id": id, "name": in.Name, "models": in.Models, "enabled": true})
}
func (s *Service) listChannels(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(r.Context(), `select id,name,base_url,models,enabled,priority,created_at,updated_at from channels order by priority,id`)
	if err != nil {
		writeError(w, 500, "internal_error", "query failed")
		return
	}
	defer rows.Close()
	data := []map[string]any{}
	for rows.Next() {
		var id, name, base string
		var models []byte
		var enabled bool
		var priority int
		var created, updated any
		rows.Scan(&id, &name, &base, &models, &enabled, &priority, &created, &updated)
		var list []string
		json.Unmarshal(models, &list)
		data = append(data, map[string]any{"id": id, "name": name, "base_url": base, "models": list, "enabled": enabled, "priority": priority, "created_at": created, "updated_at": updated})
	}
	writeJSON(w, 200, map[string]any{"data": data})
}
func (s *Service) setChannelStatus(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Enabled bool `json:"enabled"`
	}
	if decode(r, &in) != nil {
		writeError(w, 400, "invalid_request", "enabled is required")
		return
	}
	result, err := s.db.Exec(r.Context(), `update channels set enabled=$1, updated_at=now() where id=$2`, in.Enabled, r.PathValue("id"))
	if err != nil || result.RowsAffected() != 1 {
		writeError(w, 404, "not_found", "channel not found")
		return
	}
	writeJSON(w, 200, map[string]bool{"enabled": in.Enabled})
}
func (s *Service) listLogs(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(r.Context(), `select request_id,user_id,api_key_id,channel_id,model,status_code,prompt_tokens,completion_tokens,total_tokens,duration_ms,error_code,created_at from request_logs order by created_at desc limit 100`)
	if err != nil {
		writeError(w, 500, "internal_error", "query failed")
		return
	}
	defer rows.Close()
	data := []map[string]any{}
	for rows.Next() {
		var requestID, model string
		var uid, kid, cid, prompt, completion, total, errorCode, created any
		var status, duration int
		rows.Scan(&requestID, &uid, &kid, &cid, &model, &status, &prompt, &completion, &total, &duration, &errorCode, &created)
		data = append(data, map[string]any{"request_id": requestID, "user_id": uid, "api_key_id": kid, "channel_id": cid, "model": model, "status_code": status, "prompt_tokens": prompt, "completion_tokens": completion, "total_tokens": total, "duration_ms": duration, "error_code": errorCode, "created_at": created})
	}
	writeJSON(w, 200, map[string]any{"data": data})
}
