package speedtester

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParseProbeFieldMappings(t *testing.T) {
	mappings := ParseProbeFieldMappings("ip=ip,country=country_name,country_code=country_code,asn=asn,org=org")

	expected := map[string]string{
		"ip":           "ip",
		"country":      "country_name",
		"country_code": "country_code",
		"asn":          "asn",
		"org":          "org",
	}
	for key, value := range expected {
		if mappings[key] != value {
			t.Fatalf("expected mapping %s=%s, got %q", key, value, mappings[key])
		}
	}
}

func TestRunProbeExtractsMappedJSONFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Fatalf("expected GET probe request, got %s", request.Method)
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{
			"ip": "203.0.113.10",
			"country_name": "Japan",
			"country_code": "JP",
			"region": "Tokyo",
			"city": "Tokyo",
			"asn": "AS64500",
			"org": "Example Transit"
		}`))
	}))
	defer server.Close()

	result := RunProbeWithClient(server.Client(), ProbeConfig{
		URL:     server.URL,
		Method:  http.MethodGet,
		Timeout: time.Second,
		Fields:  ParseProbeFieldMappings("ip=ip,country=country_name,country_code=country_code,region=region,city=city,asn=asn,org=org"),
	})

	if result.Error != "" {
		t.Fatalf("expected no probe error, got %q", result.Error)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", result.StatusCode)
	}
	if result.Latency <= 0 {
		t.Fatalf("expected positive probe latency, got %s", result.Latency)
	}
	if result.Fields["ip"] != "203.0.113.10" {
		t.Fatalf("expected mapped ip field, got %q", result.Fields["ip"])
	}
	if result.Fields["country"] != "Japan" {
		t.Fatalf("expected mapped country field, got %q", result.Fields["country"])
	}
	if result.Fields["country_code"] != "JP" {
		t.Fatalf("expected mapped country_code field, got %q", result.Fields["country_code"])
	}
}

func TestRunProbeRecordsHTTPAndJSONErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadGateway)
		_, _ = writer.Write([]byte(`not-json`))
	}))
	defer server.Close()

	result := RunProbeWithClient(server.Client(), ProbeConfig{
		URL:     server.URL,
		Method:  http.MethodGet,
		Timeout: time.Second,
		Fields:  ParseProbeFieldMappings("ip=ip"),
	})

	if result.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected status 502, got %d", result.StatusCode)
	}
	if result.Error == "" {
		t.Fatal("expected probe error for non-200 response")
	}
}
