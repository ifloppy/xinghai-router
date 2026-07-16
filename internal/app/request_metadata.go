package app

import (
	"net"
	"net/http"
	"regexp"
	"strings"
)

type requestMetadataInfo struct {
	clientIP, forwardedFor, userAgent string
	browser, browserVersion           string
	operatingSystem, operatingSystemVersion, deviceType string
	isBot bool
}

var uaVersion = regexp.MustCompile(`(?:Chrome|Firefox|Version|Edg|OPR|CriOS|FxiOS|MSIE)[ /]([\d.]+)`)

func requestMetadataFromUA(ua string) (browser, version, operatingSystem, osVersion, device string, bot bool) {
	lower := strings.ToLower(ua)
	bot = strings.Contains(lower, "bot") || strings.Contains(lower, "crawler") || strings.Contains(lower, "spider")
	switch {
	case strings.Contains(lower, "edg/"):
		browser = "Edge"
	case strings.Contains(lower, "opr/"):
		browser = "Opera"
	case strings.Contains(lower, "chrome/") || strings.Contains(lower, "crios/"):
		browser = "Chrome"
	case strings.Contains(lower, "firefox/") || strings.Contains(lower, "fxios/"):
		browser = "Firefox"
	case strings.Contains(lower, "safari/"):
		browser = "Safari"
	case strings.Contains(lower, "msie") || strings.Contains(lower, "trident/"):
		browser = "Internet Explorer"
	default:
		browser = "Other"
	}
	if match := uaVersion.FindStringSubmatch(ua); len(match) > 1 { version = match[1] }
	switch {
	case strings.Contains(lower, "windows"):
		operatingSystem, device = "Windows", "desktop"
	case strings.Contains(lower, "android"):
		operatingSystem, device = "Android", "mobile"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad"):
		operatingSystem, device = "iOS", "mobile"
	case strings.Contains(lower, "mac os"):
		operatingSystem, device = "macOS", "desktop"
	case strings.Contains(lower, "linux"):
		operatingSystem, device = "Linux", "desktop"
	default:
		operatingSystem, device = "Other", "unknown"
	}
	return
}

func requestMetadata(r *http.Request) requestMetadataInfo {
	ua := strings.TrimSpace(r.UserAgent())
	clientIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
		if comma := strings.IndexByte(clientIP, ','); comma >= 0 { clientIP = strings.TrimSpace(clientIP[:comma]) }
	}
	if clientIP == "" {
		clientIP = r.RemoteAddr
		if host, _, err := net.SplitHostPort(clientIP); err == nil { clientIP = host }
	}
	browser, browserVersion, os, osVersion, device, bot := requestMetadataFromUA(ua)
	return requestMetadataInfo{clientIP: clientIP, forwardedFor: strings.TrimSpace(r.Header.Get("X-Forwarded-For")), userAgent: ua, browser: browser, browserVersion: browserVersion, operatingSystem: os, operatingSystemVersion: osVersion, deviceType: device, isBot: bot}
}
