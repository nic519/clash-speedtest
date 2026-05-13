package output

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/nic519/clash-speedtest/speedtester"
)

var probeFieldColumns = []string{
	"ip",
	"country",
	"country_code",
	"region",
	"city",
	"asn",
	"org",
}

// GetHeaders returns table headers based on speed mode.
// fast: ID, Name, Type, Latency, Proxy ID
// download: ID, Name, Type, Latency, Jitter, Packet Loss, Download Speed, Proxy ID
// full: ID, Name, Type, Latency, Jitter, Packet Loss, Download Speed, Upload Speed, Proxy ID
func GetHeaders(mode speedtester.SpeedMode) []string {
	if mode.IsFast() {
		return []string{
			"序号",
			"节点名称",
			"类型",
			"延迟",
		}
	}
	headers := []string{
		"序号",
		"节点名称",
		"类型",
		"延迟",
		"抖动",
		"丢包率",
		"下载速度",
	}
	if mode.UploadEnabled() {
		headers = append(headers, "上传速度")
	}
	return headers
}

func GetTSVHeaders(mode speedtester.SpeedMode) []string {
	return appendProbeHeaders(GetHeaders(mode))
}

// FormatRow formats a single result row without ANSI colors.
// Returns plain text strings using speedtester.Result's Format* methods.
func FormatRow(result *speedtester.Result, mode speedtester.SpeedMode, index int) []string {
	idStr := fmt.Sprintf("%d.", index+1)

	if mode.IsFast() {
		return []string{
			idStr,
			result.ProxyName,
			result.ProxyType,
			result.FormatLatency(),
		}
	}
	row := []string{
		idStr,
		result.ProxyName,
		result.ProxyType,
		result.FormatLatency(),
		result.FormatJitter(),
		result.FormatPacketLoss(),
		result.FormatDownloadSpeed(),
	}
	if mode.UploadEnabled() {
		row = append(row, result.FormatUploadSpeed())
	}
	return row
}

func FormatTSVRow(result *speedtester.Result, mode speedtester.SpeedMode, index int) []string {
	row := FormatRow(result, mode, index)
	return appendProbeColumns(row, result)
}

func appendProbeHeaders(headers []string) []string {
	headers = append(headers, "Probe URL", "Probe 延迟", "Probe 状态", "Probe 错误")
	for _, field := range probeFieldColumns {
		headers = append(headers, "probe."+field)
	}
	return append(headers, "节点ID")
}

func appendProbeColumns(row []string, result *speedtester.Result) []string {
	row = append(row, formatProbe(result)...)
	return append(row, ProxyID(result))
}

func formatProbe(result *speedtester.Result) []string {
	values := []string{"", "", "", ""}
	for range probeFieldColumns {
		values = append(values, "")
	}
	if result == nil || result.Probe == nil {
		return values
	}
	probe := result.Probe
	values[0] = probe.URL
	values[1] = formatDuration(probe.Latency)
	if probe.StatusCode > 0 {
		values[2] = strconv.Itoa(probe.StatusCode)
	}
	values[3] = probe.Error
	for index, field := range probeFieldColumns {
		values[4+index] = probe.Fields[field]
	}
	return values
}

func formatDuration(value time.Duration) string {
	if value <= 0 {
		return ""
	}
	return fmt.Sprintf("%dms", value.Milliseconds())
}

// SortResults sorts results based on speed mode.
// fast: latency ascending (lower is better)
// download/full: download speed descending (higher is better)
func SortResults(results []*speedtester.Result, mode speedtester.SpeedMode) []*speedtester.Result {
	if mode.IsFast() {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Latency < results[j].Latency
		})
	} else {
		sort.Slice(results, func(i, j int) bool {
			return results[i].DownloadSpeed > results[j].DownloadSpeed
		})
	}
	return results
}
