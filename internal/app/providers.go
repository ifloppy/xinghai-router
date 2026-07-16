package app

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

type modelProvider struct {
	ID, Name, Slug string
	Prefixes       []string
	Priority       int
}

func (s *Service) providers(r *http.Request) []modelProvider {
	rows, err := s.db.Query(r.Context(), `select id::text,name,slug,prefixes,priority from model_providers order by priority,name`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	data := []modelProvider{}
	for rows.Next() {
		var item modelProvider
		var prefixes []byte
		if rows.Scan(&item.ID, &item.Name, &item.Slug, &prefixes, &item.Priority) == nil {
			_ = json.Unmarshal(prefixes, &item.Prefixes)
			data = append(data, item)
		}
	}
	return data
}

func (s *Service) listProviders(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(r.Context(), `select id::text,name,slug,prefixes,priority from model_providers order by priority,name`)
	if err != nil {
		writeError(w, 500, "internal_error", "could not load providers")
		return
	}
	defer rows.Close()
	data := []modelProvider{}
	for rows.Next() {
		var item modelProvider
		var prefixes []byte
		if rows.Scan(&item.ID, &item.Name, &item.Slug, &prefixes, &item.Priority) == nil {
			_ = json.Unmarshal(prefixes, &item.Prefixes)
			data = append(data, item)
		}
	}
	writeJSON(w, 200, map[string]any{"data": data})
}

func (s *Service) saveProvider(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name, Slug string
		Prefixes   []string `json:"prefixes"`
		Priority   int      `json:"priority"`
	}
	if decode(r, &in) != nil || strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Slug) == "" || len(in.Prefixes) == 0 {
		writeError(w, 400, "invalid_request", "name, slug, and prefixes are required")
		return
	}
	for i := range in.Prefixes {
		in.Prefixes[i] = strings.ToLower(strings.TrimSpace(in.Prefixes[i]))
		if in.Prefixes[i] == "" {
			writeError(w, 400, "invalid_request", "prefixes cannot be empty")
			return
		}
	}
	if in.Priority < 0 {
		writeError(w, 400, "invalid_request", "priority cannot be negative")
		return
	}
	prefixes, _ := json.Marshal(in.Prefixes)
	var item modelProvider
	var raw []byte
	err := s.db.QueryRow(r.Context(), `insert into model_providers(name,slug,prefixes,priority,updated_at) values($1,$2,$3,$4,now()) on conflict (slug) do update set name=excluded.name,prefixes=excluded.prefixes,priority=excluded.priority,updated_at=now() returning id::text,name,slug,prefixes,priority`, strings.TrimSpace(in.Name), strings.TrimSpace(in.Slug), prefixes, in.Priority).Scan(&item.ID, &item.Name, &item.Slug, &raw, &item.Priority)
	if err != nil {
		writeError(w, 409, "conflict", "provider name or slug already exists")
		return
	}
	_ = json.Unmarshal(raw, &item.Prefixes)
	s.audit(r, "provider.saved", "model_provider", item.ID, map[string]any{"name": item.Name})
	writeJSON(w, 200, item)
}

func (s *Service) deleteProvider(w http.ResponseWriter, r *http.Request) {
	result, err := s.db.Exec(r.Context(), `delete from model_providers where id=$1`, r.PathValue("id"))
	if err != nil || result.RowsAffected() != 1 {
		writeError(w, 404, "not_found", "provider not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func providerForModel(model string, providers []modelProvider) modelProvider {
	name := strings.ToLower(model)
	matches := append([]modelProvider(nil), providers...)
	sort.SliceStable(matches, func(i, j int) bool { return matches[i].Priority < matches[j].Priority })
	for _, item := range matches {
		for _, prefix := range item.Prefixes {
			if strings.HasPrefix(name, strings.ToLower(prefix)) {
				return item
			}
		}
	}
	return modelProvider{Name: "其他", Slug: "other"}
}
