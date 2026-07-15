package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type channel struct{ id, baseURL, apiKey string }

func (s *Service) models(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(r.Context(), `select models from channels where enabled order by priority,id`)
	if err != nil {
		writeError(w, 500, "internal_error", "query failed")
		return
	}
	defer rows.Close()
	seen := map[string]bool{}
	data := []map[string]any{}
	for rows.Next() {
		var raw []byte
		if rows.Scan(&raw) != nil {
			continue
		}
		var models []string
		if json.Unmarshal(raw, &models) != nil {
			continue
		}
		for _, model := range models {
			if !seen[model] {
				seen[model] = true
				data = append(data, map[string]any{"id": model, "object": "model", "created": 0, "owned_by": "xinghai"})
			}
		}
	}
	writeJSON(w, 200, map[string]any{"object": "list", "data": data})
}

func (s *Service) chatCompletions(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	key := r.Context().Value(contextKey{}).(keyContext)
	body, err := io.ReadAll(io.LimitReader(r.Body, 2<<20))
	if err != nil {
		writeError(w, 400, "invalid_request", "could not read request")
		return
	}
	var request struct {
		Model  string `json:"model"`
		Stream bool   `json:"stream"`
	}
	if json.Unmarshal(body, &request) != nil || request.Model == "" {
		writeError(w, 400, "invalid_request", "model is required")
		return
	}
	ch, err := s.channelForModel(r, request.Model)
	if err != nil {
		s.logRequest(r, key, "", request.Model, 503, 0, 0, 0, time.Since(started), "no_channel")
		writeError(w, 503, "model_unavailable", "no enabled channel supports this model")
		return
	}
	upstreamURL := ch.baseURL + "/v1/chat/completions"
	upstreamReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upstreamURL, bytes.NewReader(body))
	if err != nil {
		s.logRequest(r, key, ch.id, request.Model, 500, 0, 0, 0, time.Since(started), "upstream_request")
		writeError(w, 500, "internal_error", "could not create upstream request")
		return
	}
	upstreamReq.Header.Set("Authorization", "Bearer "+ch.apiKey)
	upstreamReq.Header.Set("Content-Type", "application/json")
	upstreamReq.Header.Set("Accept", "application/json")
	resp, err := s.httpClient.Do(upstreamReq)
	if err != nil {
		s.logRequest(r, key, ch.id, request.Model, 502, 0, 0, 0, time.Since(started), "upstream_unreachable")
		writeError(w, 502, "upstream_error", "upstream request failed")
		return
	}
	defer resp.Body.Close()
	if request.Stream && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.streamResponse(w, resp)
		s.logRequest(r, key, ch.id, request.Model, resp.StatusCode, 0, 0, 0, time.Since(started), "")
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		writeError(w, 502, "upstream_error", "could not read upstream response")
		return
	}
	prompt, completion, total := usage(responseBody)
	s.logRequest(r, key, ch.id, request.Model, resp.StatusCode, prompt, completion, total, time.Since(started), errorCode(resp.StatusCode))
	w.Header().Set("Content-Type", contentType(resp.Header.Get("Content-Type")))
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}
func (s *Service) channelForModel(r *http.Request, model string) (channel, error) {
	var ch channel
	var encrypted string
	err := s.db.QueryRow(r.Context(), `select id,base_url,api_key from channels where enabled and models ? $1 order by priority,id limit 1`, model).Scan(&ch.id, &ch.baseURL, &encrypted)
	if err != nil {
		return ch, err
	}
	ch.apiKey, err = crypt(s.cfg.EncryptionKey, encrypted, true)
	return ch, err
}
func (s *Service) streamResponse(w http.ResponseWriter, resp *http.Response) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, 500, "internal_error", "streaming unsupported")
		return
	}
	w.Header().Set("Content-Type", contentType(resp.Header.Get("Content-Type")))
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(resp.StatusCode)
	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return
			}
			flusher.Flush()
		}
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}
	}
}
func (s *Service) logRequest(r *http.Request, key keyContext, channelID, model string, status, prompt, completion, total int, d time.Duration, errorCode string) {
	id, _ := randomID()
	_, _ = s.db.Exec(r.Context(), `insert into request_logs(id,request_id,user_id,api_key_id,channel_id,model,status_code,prompt_tokens,completion_tokens,total_tokens,duration_ms,error_code) values($1,$2,$3,$4,nullif($5,'')::uuid,$6,$7,$8,$9,$10,$11,nullif($12,''))`, id, requestID(r.Context()), key.userID, key.keyID, channelID, model, status, prompt, completion, total, d.Milliseconds(), errorCode)
	_, _ = s.db.Exec(r.Context(), `update api_keys set last_used_at=now() where id=$1`, key.keyID)
}
func usage(body []byte) (int, int, int) {
	var v struct {
		Usage struct {
			Prompt     int `json:"prompt_tokens"`
			Completion int `json:"completion_tokens"`
			Total      int `json:"total_tokens"`
		} `json:"usage"`
	}
	if json.Unmarshal(body, &v) != nil {
		return 0, 0, 0
	}
	return v.Usage.Prompt, v.Usage.Completion, v.Usage.Total
}
func errorCode(status int) string {
	if status >= 400 {
		return "upstream_" + http.StatusText(status)
	}
	return ""
}
func contentType(value string) string {
	if strings.HasPrefix(value, "application/json") {
		return "application/json"
	}
	if strings.HasPrefix(value, "text/event-stream") {
		return "text/event-stream"
	}
	return "application/json"
}
