package app

import (
	"testing"
	"time"
)

func TestRankingDuration(t *testing.T) {
	for period, expected := range map[string]time.Duration{"today": 24 * time.Hour, "week": 7 * 24 * time.Hour, "month": 30 * 24 * time.Hour, "year": 365 * 24 * time.Hour} {
		actual, ok := rankingDuration(period)
		if !ok || actual != expected {
			t.Fatalf("rankingDuration(%q) = %v, %v", period, actual, ok)
		}
	}
	if _, ok := rankingDuration("all"); ok {
		t.Fatal("unsupported period was accepted")
	}
}
func TestModelVendor(t *testing.T) {
	for model, expected := range map[string]string{"kimi-k3": "其他", "claude-3": "Anthropic", "gemini-2.5": "Google", "qwen-max": "Alibaba", "custom": "其他"} {
		if actual := modelVendor(model); actual != expected {
			t.Errorf("modelVendor(%q) = %q, want %q", model, actual, expected)
		}
	}
}
func TestGrowthPercent(t *testing.T) {
	if growthPercent(150, 100) != 50 || growthPercent(10, 0) != 100 || growthPercent(0, 0) != 0 {
		t.Fatal("unexpected growth calculation")
	}
}
