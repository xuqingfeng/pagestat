package trace

import (
	"testing"
	"time"
)

func TestTrace(t *testing.T) {

	l := make(map[string]time.Duration)
	l["dns_lookup"] = 0
	l["tcp_connection"] = 0
	l["tls_handshake"] = 0
	l["server_processing"] = 0
	err := Trace("GET", "https://example.com", []string{}, "", 0, l)
	if err != nil {
		t.Errorf("E! trace fail %v", err)
	}
	t.Logf("I! latency %+v", l)
}
