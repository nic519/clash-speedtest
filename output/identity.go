package output

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/nic519/clash-speedtest/speedtester"
)

// ProxyID returns a stable, non-secret identifier for a proxy. It hashes the
// canonical connection identity instead of exposing credentials in TSV output.
func ProxyID(result *speedtester.Result) string {
	if result == nil {
		return ""
	}

	parts := []string{
		"type=" + valueString(result.ProxyConfig["type"]),
		"server=" + valueString(result.ProxyConfig["server"]),
		"port=" + valueString(result.ProxyConfig["port"]),
	}

	for _, key := range []string{
		"network",
		"cipher",
		"uuid",
		"password",
		"username",
		"alterId",
		"sni",
		"servername",
		"ws-opts",
		"grpc-opts",
		"reality-opts",
	} {
		if value, ok := result.ProxyConfig[key]; ok {
			parts = append(parts, key+"="+valueString(value))
		}
	}

	if len(result.ProxyConfig) == 0 {
		parts = append(parts, "name="+result.ProxyName, "proxy_type="+result.ProxyType)
	}

	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(sum[:])[:16]
}

func valueString(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			parts = append(parts, key+"="+valueString(typed[key]))
		}
		return "{" + strings.Join(parts, ",") + "}"
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			parts = append(parts, valueString(item))
		}
		return "[" + strings.Join(parts, ",") + "]"
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", typed))
	}
}
