package app

import (
	"net/http"
	"net/url"
	"strings"
)

func (s *Service) siteSettings(w http.ResponseWriter, r *http.Request) {
	var name, iconURL string
	if err := s.db.QueryRow(r.Context(), `select name,icon_url from site_settings where id=true`).Scan(&name, &iconURL); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not load site settings")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"name": name, "icon_url": iconURL})
}

func (s *Service) updateSiteSettings(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name    string `json:"name"`
		IconURL string `json:"icon_url"`
	}
	if decode(r, &in) != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid site settings")
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.IconURL = strings.TrimSpace(in.IconURL)
	if in.Name == "" || len([]rune(in.Name)) > 100 {
		writeError(w, http.StatusBadRequest, "invalid_request", "site name must contain 1 to 100 characters")
		return
	}
	if in.IconURL != "" {
		u, err := url.Parse(in.IconURL)
		if err != nil || u.Host == "" || (u.Scheme != "https" && !(u.Scheme == "http" && isLoopbackHost(u.Hostname()))) {
			writeError(w, http.StatusBadRequest, "invalid_request", "icon_url must use HTTPS, except for loopback HTTP URLs")
			return
		}
	}
	if _, err := s.db.Exec(r.Context(), `update site_settings set name=$1,icon_url=$2,updated_at=now() where id=true`, in.Name, in.IconURL); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not save site settings")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"name": in.Name, "icon_url": in.IconURL})
}
