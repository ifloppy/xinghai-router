package app

import (
	"net/http"
	"net/url"
	"strings"
)

func (s *Service) siteSettings(w http.ResponseWriter, r *http.Request) {
	var name, iconURL string
	var autoDisableFailedChannels bool
	if err := s.db.QueryRow(r.Context(), `select name,icon_url,auto_disable_failed_channels from site_settings where id=true`).Scan(&name, &iconURL, &autoDisableFailedChannels); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not load site settings")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"name": name, "icon_url": iconURL, "auto_disable_failed_channels": autoDisableFailedChannels})
}

func (s *Service) updateSiteSettings(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name                      string `json:"name"`
		IconURL                   string `json:"icon_url"`
		AutoDisableFailedChannels *bool  `json:"auto_disable_failed_channels"`
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
	if _, err := s.db.Exec(r.Context(), `update site_settings set name=$1,icon_url=$2,auto_disable_failed_channels=coalesce($3,auto_disable_failed_channels),updated_at=now() where id=true`, in.Name, in.IconURL, in.AutoDisableFailedChannels); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not save site settings")
		return
	}
	var autoDisableFailedChannels bool
	if err := s.db.QueryRow(r.Context(), `select auto_disable_failed_channels from site_settings where id=true`).Scan(&autoDisableFailedChannels); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not load site settings")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"name": in.Name, "icon_url": in.IconURL, "auto_disable_failed_channels": autoDisableFailedChannels})
}
