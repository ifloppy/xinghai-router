package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOptionalAccountAllowsAnonymousRequest(t *testing.T) {
	s := &Service{}
	called := false
	handler := s.optionalAccount(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if account := accountFromContext(r); account.userID != "" {
			t.Fatalf("anonymous request has user ID %q", account.userID)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/model-catalog", nil))

	if !called {
		t.Fatal("model catalog handler was not called")
	}
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
}

func TestValidatePasswordChange(t *testing.T) {
	if msg := validatePasswordChange("old-password", "new-password"); msg != "" {
		t.Fatalf("expected valid change, got %q", msg)
	}
	cases := []struct {
		current, next string
	}{
		{"", "new-password"},
		{"old-password", ""},
		{"old-password", "short"},
		{"old-password", strings.Repeat("a", 73)},
		{"same-password", "same-password"},
	}
	for _, tc := range cases {
		if msg := validatePasswordChange(tc.current, tc.next); msg == "" {
			t.Fatalf("expected rejection for current=%q next=%q", tc.current, tc.next)
		}
	}
}

type sequenceLimiter struct {
	remaining int
}

func (l *sequenceLimiter) allow(key string) bool {
	if l.remaining <= 0 {
		return false
	}
	l.remaining--
	return true
}

func (l *sequenceLimiter) close() {}

func TestLoginRateLimitBeforeDatabase(t *testing.T) {
	// remaining=0 forces 429 before captcha/DB.
	s := &Service{limiter: &sequenceLimiter{remaining: 0}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"password1"}`))
	req.RemoteAddr = "203.0.113.10:12345"
	s.login(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{}`))
	(&Service{}).login(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("empty login status = %d", rec.Code)
	}
}

func TestDummyPasswordHashSpendsTime(t *testing.T) {
	if dummyPasswordHash == "" {
		t.Fatal("dummy password hash must be initialized")
	}
	if passwordMatches(dummyPasswordHash, "password1") {
		t.Fatal("dummy hash must not match arbitrary passwords")
	}
}

func TestChangeAccountPasswordRejectsInvalidBodyBeforeDatabaseAccess(t *testing.T) {
	for _, body := range []string{
		`{}`,
		`{"current_password":"old-password"}`,
		`{"new_password":"new-password"}`,
		`{"current_password":"old-password","new_password":"short"}`,
		`{"current_password":"same-pass","new_password":"same-pass"}`,
	} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/account/password", strings.NewReader(body))
		request = request.WithContext(context.WithValue(request.Context(), accountContextKey{}, accountContext{userID: "1"}))
		(&Service{}).changeAccountPassword(recorder, request)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("body %s status = %d, want %d", body, recorder.Code, http.StatusBadRequest)
		}
	}
}


