package app

import "testing"

func TestRequestMetadataFromUA(t *testing.T) {
	browser, version, os, _, device, bot := requestMetadataFromUA("Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 Version/17.4 Mobile/15E148 Safari/604.1")
	if browser != "Safari" || version != "17.4" || os != "iOS" || device != "mobile" || bot {
		t.Fatalf("unexpected metadata: browser=%q version=%q os=%q device=%q bot=%v", browser, version, os, device, bot)
	}
	_, _, _, _, _, bot = requestMetadataFromUA("ExampleBot/1.0")
	if !bot {
		t.Fatal("expected bot user agent")
	}
}
