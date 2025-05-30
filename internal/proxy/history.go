package proxy

import (
	"encoding/json"
	"sync"
	"time"
)

// RequestRecord represents a single HTTP request and response through the proxy
type RequestRecord struct {
	ID              string            `json:"id"`
	Timestamp       time.Time         `json:"timestamp"`
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	RequestHeaders  map[string]string `json:"request_headers"`
	RequestBody     string            `json:"request_body,omitempty"`
	ResponseStatus  int               `json:"response_status"`
	ResponseHeaders map[string]string `json:"response_headers"`
	ResponseBody    string            `json:"response_body,omitempty"`

	// Timing metrics
	ProxyStartTime    time.Time `json:"proxy_start_time"`
	UpstreamStartTime time.Time `json:"upstream_start_time"`
	UpstreamEndTime   time.Time `json:"upstream_end_time"`
	ProxyEndTime      time.Time `json:"proxy_end_time"`

	// Calculated metrics
	ProxyOverheadMs   int64 `json:"proxy_overhead_ms"`   // Time spent in proxy logic
	UpstreamLatencyMs int64 `json:"upstream_latency_ms"` // Time waiting for upstream
	TotalDurationMs   int64 `json:"total_duration_ms"`   // Total time from client perspective

	// Size metrics
	RequestSize  int64 `json:"request_size"`
	ResponseSize int64 `json:"response_size"`

	// Status
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// RequestHistory manages the collection of request records
type RequestHistory struct {
	records []RequestRecord
	mutex   sync.RWMutex
	maxSize int
}

// NewRequestHistory creates a new request history with the specified maximum size
func NewRequestHistory(maxSize int) *RequestHistory {
	return &RequestHistory{
		records: make([]RequestRecord, 0),
		maxSize: maxSize,
	}
}

// AddRecord adds a new request record to the history
func (h *RequestHistory) AddRecord(record RequestRecord) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Calculate metrics
	record.TotalDurationMs = record.ProxyEndTime.Sub(record.ProxyStartTime).Milliseconds()
	record.UpstreamLatencyMs = record.UpstreamEndTime.Sub(record.UpstreamStartTime).Milliseconds()
	record.ProxyOverheadMs = record.TotalDurationMs - record.UpstreamLatencyMs

	// Add to beginning of slice (most recent first)
	h.records = append([]RequestRecord{record}, h.records...)

	// Trim to max size
	if len(h.records) > h.maxSize {
		h.records = h.records[:h.maxSize]
	}
}

// GetRecords returns all records (most recent first)
func (h *RequestHistory) GetRecords() []RequestRecord {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]RequestRecord, len(h.records))
	copy(result, h.records)
	return result
}

// GetRecordsJSON returns all records as JSON
func (h *RequestHistory) GetRecordsJSON() ([]byte, error) {
	records := h.GetRecords()
	return json.Marshal(map[string]interface{}{
		"records": records,
		"total":   len(records),
	})
}

// Clear removes all records
func (h *RequestHistory) Clear() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.records = h.records[:0]
}

// GetStats returns aggregated statistics
func (h *RequestHistory) GetStats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if len(h.records) == 0 {
		return map[string]interface{}{
			"total_requests": 0,
		}
	}

	var totalDuration, totalUpstreamLatency, totalProxyOverhead int64
	var totalRequestSize, totalResponseSize int64
	var successCount, errorCount int
	statusCounts := make(map[int]int)
	methodCounts := make(map[string]int)

	for _, record := range h.records {
		totalDuration += record.TotalDurationMs
		totalUpstreamLatency += record.UpstreamLatencyMs
		totalProxyOverhead += record.ProxyOverheadMs
		totalRequestSize += record.RequestSize
		totalResponseSize += record.ResponseSize

		if record.Success {
			successCount++
		} else {
			errorCount++
		}

		statusCounts[record.ResponseStatus]++
		methodCounts[record.Method]++
	}

	count := len(h.records)
	return map[string]interface{}{
		"total_requests":          count,
		"success_count":           successCount,
		"error_count":             errorCount,
		"avg_duration_ms":         totalDuration / int64(count),
		"avg_upstream_latency_ms": totalUpstreamLatency / int64(count),
		"avg_proxy_overhead_ms":   totalProxyOverhead / int64(count),
		"total_request_size":      totalRequestSize,
		"total_response_size":     totalResponseSize,
		"status_codes":            statusCounts,
		"methods":                 methodCounts,
	}
}
