package speedtester

import (
	"strings"
	"testing"
)

func TestResolveServerTarget(t *testing.T) {
	t.Run("download server without path", func(t *testing.T) {
		target, err := resolveServerTarget("https://example.com")
		if err != nil {
			t.Fatalf("resolveServerTarget failed: %v", err)
		}
		if target.mode != serverModeDownloadServer {
			t.Fatalf("expected download server mode, got %v", target.mode)
		}
		if target.baseURL != "https://example.com" {
			t.Fatalf("expected baseURL to be trimmed, got %q", target.baseURL)
		}
	})

	t.Run("download server with trailing slash", func(t *testing.T) {
		target, err := resolveServerTarget("https://example.com/")
		if err != nil {
			t.Fatalf("resolveServerTarget failed: %v", err)
		}
		if target.mode != serverModeDownloadServer {
			t.Fatalf("expected download server mode, got %v", target.mode)
		}
		if target.baseURL != "https://example.com" {
			t.Fatalf("expected baseURL to be trimmed, got %q", target.baseURL)
		}
	})

	t.Run("direct download with path", func(t *testing.T) {
		target, err := resolveServerTarget("https://example.com/file.bin")
		if err != nil {
			t.Fatalf("resolveServerTarget failed: %v", err)
		}
		if target.mode != serverModeDirectDownload {
			t.Fatalf("expected direct download mode, got %v", target.mode)
		}
		if target.downloadURL != "https://example.com/file.bin" {
			t.Fatalf("expected downloadURL to be preserved, got %q", target.downloadURL)
		}
	})

	t.Run("direct download with query", func(t *testing.T) {
		target, err := resolveServerTarget("https://example.com?bytes=1024")
		if err != nil {
			t.Fatalf("resolveServerTarget failed: %v", err)
		}
		if target.mode != serverModeDirectDownload {
			t.Fatalf("expected direct download mode, got %v", target.mode)
		}
		if target.downloadURL != "https://example.com?bytes=1024" {
			t.Fatalf("expected downloadURL to be preserved, got %q", target.downloadURL)
		}
	})

	t.Run("invalid scheme", func(t *testing.T) {
		_, err := resolveServerTarget("ftp://example.com")
		if err == nil {
			t.Fatal("expected error for invalid scheme")
		}
		if !strings.Contains(err.Error(), "ftp://example.com") {
			t.Fatalf("expected error to include server url context, got %v", err)
		}
	})
}

func TestNewDisablesUploadForDirectURL(t *testing.T) {
	st, err := New(&Config{
		ServerURL:    "https://example.com/file.bin",
		UploadSize:   10,
		DownloadSize: 10,
		Concurrent:   1,
		Mode:         SpeedModeFull,
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if st.Mode() != SpeedModeDownload {
		t.Fatalf("expected mode to downgrade to %s for direct download url, got %s", SpeedModeDownload, st.Mode())
	}
}

func TestNewEnablesUploadForDownloadServer(t *testing.T) {
	st, err := New(&Config{
		ServerURL:    "https://example.com",
		UploadSize:   10,
		DownloadSize: 10,
		Concurrent:   1,
		Mode:         SpeedModeFull,
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if st.Mode() != SpeedModeFull {
		t.Fatalf("expected mode to remain %s for download server url, got %s", SpeedModeFull, st.Mode())
	}
}

func TestNewDefaultsLatencyURLToDownloadTarget(t *testing.T) {
	t.Run("direct download url", func(t *testing.T) {
		st, err := New(&Config{
			ServerURL:    "https://example.com/file.bin",
			DownloadSize: 10,
			Concurrent:   1,
			Mode:         SpeedModeFast,
		})
		if err != nil {
			t.Fatalf("New failed: %v", err)
		}
		if st.latencyURL != "https://example.com/file.bin" {
			t.Fatalf("expected latencyURL to default to direct download url, got %q", st.latencyURL)
		}
	})

	t.Run("download server url", func(t *testing.T) {
		st, err := New(&Config{
			ServerURL:    "https://example.com",
			DownloadSize: 10,
			Concurrent:   1,
			Mode:         SpeedModeFast,
		})
		if err != nil {
			t.Fatalf("New failed: %v", err)
		}
		if st.latencyURL != "https://example.com" {
			t.Fatalf("expected latencyURL to default to server base url, got %q", st.latencyURL)
		}
	})
}

func TestNewUsesConfiguredLatencyURL(t *testing.T) {
	st, err := New(&Config{
		ServerURL:    "https://example.com/file.bin",
		LatencyURL:   "https://youtube.com/generate_204",
		DownloadSize: 10,
		Concurrent:   1,
		Mode:         SpeedModeFast,
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if st.latencyURL != "https://youtube.com/generate_204" {
		t.Fatalf("expected latencyURL to use configured value, got %q", st.latencyURL)
	}
	if st.downloadURL != "https://example.com/file.bin" {
		t.Fatalf("expected downloadURL to remain unchanged, got %q", st.downloadURL)
	}
}

func TestNewRejectsInvalidLatencyURL(t *testing.T) {
	_, err := New(&Config{
		ServerURL:    "https://example.com/file.bin",
		LatencyURL:   "ftp://youtube.com/generate_204",
		DownloadSize: 10,
		Concurrent:   1,
		Mode:         SpeedModeFast,
	})
	if err == nil {
		t.Fatal("expected error for invalid latency url")
	}
	if !strings.Contains(err.Error(), "latency url") {
		t.Fatalf("expected error to include latency url context, got %v", err)
	}
}
