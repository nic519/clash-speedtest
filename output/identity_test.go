package output

import (
	"testing"

	"github.com/nic519/clash-speedtest/speedtester"
)

func TestProxyIDUsesConnectionIdentity(t *testing.T) {
	first := &speedtester.Result{
		ProxyName: "Renamable",
		ProxyType: "Trojan",
		ProxyConfig: map[string]any{
			"type":     "trojan",
			"server":   "example.com",
			"port":     443,
			"password": "secret",
		},
	}
	renamed := &speedtester.Result{
		ProxyName: "Different Display Name",
		ProxyType: "Trojan",
		ProxyConfig: map[string]any{
			"type":     "trojan",
			"server":   "example.com",
			"port":     443,
			"password": "secret",
		},
	}
	differentPort := &speedtester.Result{
		ProxyName: "Renamable",
		ProxyType: "Trojan",
		ProxyConfig: map[string]any{
			"type":     "trojan",
			"server":   "example.com",
			"port":     8443,
			"password": "secret",
		},
	}

	if ProxyID(first) == "" {
		t.Fatal("expected non-empty proxy id")
	}
	if ProxyID(first) != ProxyID(renamed) {
		t.Fatal("expected display name changes to keep the same proxy id")
	}
	if ProxyID(first) == ProxyID(differentPort) {
		t.Fatal("expected connection changes to produce a different proxy id")
	}
}
