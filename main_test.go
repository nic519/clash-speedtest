package main

import "testing"

func TestResolveVersionPrefersInjectedVersion(t *testing.T) {
	got := resolveVersion("0.1.4", "v0.1.5")
	if got != "0.1.4" {
		t.Fatalf("expected injected version, got %q", got)
	}
}

func TestResolveVersionFallsBackToBuildInfoVersion(t *testing.T) {
	got := resolveVersion("dev", "v0.1.4")
	if got != "0.1.4" {
		t.Fatalf("expected build info version, got %q", got)
	}
}

func TestResolveVersionKeepsDevWithoutBuildInfo(t *testing.T) {
	got := resolveVersion("dev", "(devel)")
	if got != "dev" {
		t.Fatalf("expected dev fallback, got %q", got)
	}
}
