package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseCreditAmount(t *testing.T) {
	tests := []struct {
		value string
		want  float64
		ok    bool
	}{
		{"", 0, true},
		{"10", 10, true},
		{"10.5", 10.5, true},
		{"0.01", 0.01, true},
		{"-1", 0, false},
		{"1.001", 0, false},
		{"abc", 0, false},
	}
	for _, tt := range tests {
		got, ok := parseCreditAmount(tt.value)
		if ok != tt.ok || (ok && got != tt.want) {
			t.Fatalf("parseCreditAmount(%q) = (%v, %v), want (%v, %v)", tt.value, got, ok, tt.want, tt.ok)
		}
	}
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		cents int64
		want  string
	}{
		{0, "0.00"},
		{100, "1.00"},
		{1050, "10.50"},
		{1, "0.01"},
		{10000000, "100000.00"},
	}
	for _, tt := range tests {
		if got := formatAmount(tt.cents); got != tt.want {
			t.Fatalf("formatAmount(%d) = %q, want %q", tt.cents, got, tt.want)
		}
	}
}

func TestFormatCredit(t *testing.T) {
	if got := formatCredit(10.5); got != "10.5" {
		t.Fatalf("formatCredit(10.5) = %q, want %q", got, "10.5")
	}
	if got := formatCredit(0); got != "0" {
		t.Fatalf("formatCredit(0) = %q, want %q", got, "0")
	}
}

func TestReadSubscriptionPlanInputRejectsOutOfRange(t *testing.T) {
	cases := []string{
		`{"name":"Plan","price":"100000.01","billing_period":"month","credit_amount":"0"}`,
		`{"name":"Plan","price":"1","billing_period":"month","credit_amount":"1000000.01"}`,
		`{"name":"Plan","price":"1","billing_period":"month","credit_amount":"0","sort_order":10001}`,
		`{"name":"Plan","price":"1","billing_period":"month","credit_amount":"0","max_requests_per_period":-1}`,
		`{"name":"Plan","price":"1","billing_period":"month","credit_amount":"0","description":"` + strings.Repeat("d", 2001) + `"}`,
	}
	for _, body := range cases {
		req := httptest.NewRequest(http.MethodPost, "/admin/subscription-plans", strings.NewReader(body))
		if _, err := readSubscriptionPlanInput(req, &Service{}, ""); err == nil {
			t.Fatalf("expected rejection for body %s", body)
		}
	}
}

func TestValidPlanQuotaLimit(t *testing.T) {
	var ok int64 = maxPlanQuotaLimit
	var over int64 = maxPlanQuotaLimit + 1
	var neg int64 = -1
	if !validPlanQuotaLimit(nil) || !validPlanQuotaLimit(&ok) || validPlanQuotaLimit(&over) || validPlanQuotaLimit(&neg) {
		t.Fatal("plan quota limit bounds unexpected")
	}
}
