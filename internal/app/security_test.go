package app

import (
	"strings"
	"testing"
	"time"
)

func TestCryptRoundTrip(t *testing.T) {
	encoded, err := crypt("a sufficiently long test encryption key", "provider-secret", false)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(encoded, "provider-secret") {
		t.Fatal("ciphertext contains plaintext")
	}
	decoded, err := crypt("a sufficiently long test encryption key", encoded, true)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != "provider-secret" {
		t.Fatalf("got %q", decoded)
	}
}

func TestLimiter(t *testing.T) {
	l := newLimiter(2)
	if !l.allow("key") || !l.allow("key") || l.allow("key") {
		t.Fatal("unexpected rate-limit result")
	}
	l.entries["key"] = rateWindow{start: time.Now().Add(-time.Minute)}
	if !l.allow("key") {
		t.Fatal("window did not reset")
	}
}

func TestUsage(t *testing.T) {
	prompt, completion, total := usage([]byte(`{"usage":{"prompt_tokens":2,"completion_tokens":3,"total_tokens":5}}`))
	if prompt != 2 || completion != 3 || total != 5 {
		t.Fatalf("got %d, %d, %d", prompt, completion, total)
	}
}
