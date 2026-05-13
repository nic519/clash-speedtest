package speedtester

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ProbeConfig struct {
	URL     string
	Method  string
	Timeout time.Duration
	Fields  map[string]string
}

type ProbeResult struct {
	URL        string            `json:"url"`
	Method     string            `json:"method"`
	Latency    time.Duration     `json:"latency"`
	StatusCode int               `json:"status_code"`
	Error      string            `json:"error"`
	Fields     map[string]string `json:"fields"`
}

func ParseProbeFieldMappings(raw string) map[string]string {
	mappings := make(map[string]string)
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue
		}
		mappings[key] = value
	}
	return mappings
}

func RunProbeWithClient(client *http.Client, config ProbeConfig) *ProbeResult {
	method := strings.TrimSpace(config.Method)
	if method == "" {
		method = http.MethodGet
	}
	result := &ProbeResult{
		URL:    strings.TrimSpace(config.URL),
		Method: method,
		Fields: make(map[string]string),
	}
	if result.URL == "" {
		result.Error = "probe url is empty"
		return result
	}
	if client == nil {
		client = http.DefaultClient
	}
	probeClient := *client
	if config.Timeout > 0 {
		probeClient.Timeout = config.Timeout
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if config.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}

	request, err := http.NewRequestWithContext(ctx, method, result.URL, nil)
	if err != nil {
		result.Error = fmt.Sprintf("create probe request failed: %v", err)
		return result
	}
	request.Header.Set("User-Agent", "clash-speedtest/1.0")

	start := time.Now()
	response, err := probeClient.Do(request)
	result.Latency = time.Since(start)
	if err != nil {
		result.Error = fmt.Sprintf("probe request failed: %v", err)
		return result
	}
	defer response.Body.Close()
	result.StatusCode = response.StatusCode

	body, err := io.ReadAll(response.Body)
	if err != nil {
		result.Error = fmt.Sprintf("read probe response failed: %v", err)
		return result
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		result.Error = fmt.Sprintf("probe response returned %s", response.Status)
		return result
	}

	fields, err := extractProbeFields(body, config.Fields)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	result.Fields = fields
	return result
}

func extractProbeFields(body []byte, mappings map[string]string) (map[string]string, error) {
	fields := make(map[string]string)
	if len(mappings) == 0 {
		return fields, nil
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		return fields, fmt.Errorf("parse probe json failed: %w", err)
	}

	for outputKey, path := range mappings {
		if value, ok := lookupJSONPath(payload, path); ok {
			fields[outputKey] = stringifyProbeValue(value)
		}
	}
	return fields, nil
}

func lookupJSONPath(value any, path string) (any, bool) {
	current := value
	for _, part := range strings.Split(path, ".") {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, false
		}
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = object[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func stringifyProbeValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(typed)
	case float64:
		if typed == float64(int64(typed)) {
			return strconv.FormatInt(int64(typed), 10)
		}
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", typed))
	}
}
