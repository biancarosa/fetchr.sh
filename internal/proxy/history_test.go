//go:build unit

package proxy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRequestHistory(t *testing.T) {
	maxSize := 5
	history := NewRequestHistory(maxSize)

	if history == nil {
		t.Fatal("NewRequestHistory returned nil")
	}

	if history.maxSize != maxSize {
		t.Errorf("Expected maxSize %d, got %d", maxSize, history.maxSize)
	}

	if len(history.records) != 0 {
		t.Errorf("Expected empty records, got length %d", len(history.records))
	}
}

func TestAddRecord(t *testing.T) {
	history := NewRequestHistory(3)

	// Create test records
	record1 := RequestRecord{
		ID:                "1",
		Timestamp:         time.Now(),
		Method:            "GET",
		URL:               "http://example.com",
		ProxyStartTime:    time.Now(),
		UpstreamStartTime: time.Now().Add(time.Millisecond),
		UpstreamEndTime:   time.Now().Add(2 * time.Millisecond),
		ProxyEndTime:      time.Now().Add(3 * time.Millisecond),
		Success:           true,
	}

	record2 := RequestRecord{
		ID:                "2",
		Timestamp:         time.Now(),
		Method:            "POST",
		URL:               "http://example.com/api",
		ProxyStartTime:    time.Now(),
		UpstreamStartTime: time.Now().Add(time.Millisecond),
		UpstreamEndTime:   time.Now().Add(2 * time.Millisecond),
		ProxyEndTime:      time.Now().Add(3 * time.Millisecond),
		Success:           true,
	}

	// Add first record
	history.AddRecord(record1)
	records := history.GetRecords()

	if len(records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(records))
	}

	if records[0].ID != "1" {
		t.Errorf("Expected first record ID '1', got '%s'", records[0].ID)
	}

	// Add second record
	history.AddRecord(record2)
	records = history.GetRecords()

	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}

	// Verify most recent first
	if records[0].ID != "2" {
		t.Errorf("Expected most recent record ID '2', got '%s'", records[0].ID)
	}

	if records[1].ID != "1" {
		t.Errorf("Expected second record ID '1', got '%s'", records[1].ID)
	}
}

func TestMaxSize(t *testing.T) {
	maxSize := 2
	history := NewRequestHistory(maxSize)

	// Add records exceeding max size
	for i := 1; i <= 4; i++ {
		record := RequestRecord{
			ID:                string(rune('0' + i)),
			Timestamp:         time.Now(),
			Method:            "GET",
			URL:               "http://example.com",
			ProxyStartTime:    time.Now(),
			UpstreamStartTime: time.Now().Add(time.Millisecond),
			UpstreamEndTime:   time.Now().Add(2 * time.Millisecond),
			ProxyEndTime:      time.Now().Add(3 * time.Millisecond),
			Success:           true,
		}
		history.AddRecord(record)
	}

	records := history.GetRecords()

	// Should only keep the last 2 records
	if len(records) != maxSize {
		t.Errorf("Expected %d records, got %d", maxSize, len(records))
	}

	// Should have the most recent ones (4 and 3)
	if records[0].ID != "4" {
		t.Errorf("Expected most recent record ID '4', got '%s'", records[0].ID)
	}

	if records[1].ID != "3" {
		t.Errorf("Expected second record ID '3', got '%s'", records[1].ID)
	}
}

func TestCalculateMetrics(t *testing.T) {
	history := NewRequestHistory(10)

	// Create a record with specific timing
	now := time.Now()
	record := RequestRecord{
		ID:                "test",
		Timestamp:         now,
		Method:            "GET",
		URL:               "http://example.com",
		ProxyStartTime:    now,
		UpstreamStartTime: now.Add(5 * time.Millisecond),
		UpstreamEndTime:   now.Add(15 * time.Millisecond),
		ProxyEndTime:      now.Add(20 * time.Millisecond),
		Success:           true,
	}

	history.AddRecord(record)
	records := history.GetRecords()

	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}

	r := records[0]

	// Total duration should be 20000µs (20ms)
	if r.TotalDurationUs != 20000 {
		t.Errorf("Expected TotalDurationUs 20000, got %d", r.TotalDurationUs)
	}

	// Upstream latency should be 10000µs (10ms)
	if r.UpstreamLatencyUs != 10000 {
		t.Errorf("Expected UpstreamLatencyUs 10000, got %d", r.UpstreamLatencyUs)
	}

	// Proxy overhead should be 10000µs (10ms = 20ms - 10ms)
	if r.ProxyOverheadUs != 10000 {
		t.Errorf("Expected ProxyOverheadUs 10000, got %d", r.ProxyOverheadUs)
	}
}

func TestClear(t *testing.T) {
	history := NewRequestHistory(10)

	// Add some records
	for i := 1; i <= 3; i++ {
		record := RequestRecord{
			ID:        string(rune('0' + i)),
			Timestamp: time.Now(),
			Method:    "GET",
			URL:       "http://example.com",
			Success:   true,
		}
		history.AddRecord(record)
	}

	if len(history.GetRecords()) != 3 {
		t.Errorf("Expected 3 records before clear, got %d", len(history.GetRecords()))
	}

	// Clear history
	history.Clear()

	if len(history.GetRecords()) != 0 {
		t.Errorf("Expected 0 records after clear, got %d", len(history.GetRecords()))
	}
}

func TestGetStats(t *testing.T) {
	history := NewRequestHistory(10)

	// Test empty stats
	stats := history.GetStats()
	if stats["total_requests"] != 0 {
		t.Errorf("Expected 0 total_requests for empty history, got %v", stats["total_requests"])
	}

	// Add test records
	now := time.Now()

	// Successful GET request
	history.AddRecord(RequestRecord{
		ID:                "1",
		Method:            "GET",
		ResponseStatus:    200,
		ProxyStartTime:    now,
		UpstreamStartTime: now.Add(5 * time.Millisecond),
		UpstreamEndTime:   now.Add(15 * time.Millisecond),
		ProxyEndTime:      now.Add(20 * time.Millisecond),
		RequestSize:       100,
		ResponseSize:      500,
		Success:           true,
	})

	// Failed POST request
	history.AddRecord(RequestRecord{
		ID:                "2",
		Method:            "POST",
		ResponseStatus:    500,
		ProxyStartTime:    now,
		UpstreamStartTime: now.Add(2 * time.Millisecond),
		UpstreamEndTime:   now.Add(12 * time.Millisecond),
		ProxyEndTime:      now.Add(15 * time.Millisecond),
		RequestSize:       200,
		ResponseSize:      100,
		Success:           false,
	})

	stats = history.GetStats()

	if stats["total_requests"] != 2 {
		t.Errorf("Expected 2 total_requests, got %v", stats["total_requests"])
	}

	if stats["success_count"] != 1 {
		t.Errorf("Expected 1 success_count, got %v", stats["success_count"])
	}

	if stats["error_count"] != 1 {
		t.Errorf("Expected 1 error_count, got %v", stats["error_count"])
	}

	// Check averages: (20000+15000)/2 = 17500µs
	if stats["avg_duration_us"] != int64(17500) {
		t.Errorf("Expected avg_duration_us 17500, got %v", stats["avg_duration_us"])
	}

	// Check status codes
	statusCodes := stats["status_codes"].(map[int]int)
	if statusCodes[200] != 1 {
		t.Errorf("Expected 1 request with status 200, got %d", statusCodes[200])
	}
	if statusCodes[500] != 1 {
		t.Errorf("Expected 1 request with status 500, got %d", statusCodes[500])
	}

	// Check methods
	methods := stats["methods"].(map[string]int)
	if methods["GET"] != 1 {
		t.Errorf("Expected 1 GET request, got %d", methods["GET"])
	}
	if methods["POST"] != 1 {
		t.Errorf("Expected 1 POST request, got %d", methods["POST"])
	}
}

func TestProxyOverheadCalculation(t *testing.T) {
	history := NewRequestHistory(10)

	// Create a test record with known timing values
	now := time.Now()
	record := RequestRecord{
		ID:                "test-overhead",
		Method:            "GET",
		URL:               "http://example.com",
		ProxyStartTime:    now,                            // Start: 0µs
		UpstreamStartTime: now.Add(5 * time.Millisecond),  // Proxy processing: 5000µs
		UpstreamEndTime:   now.Add(15 * time.Millisecond), // Upstream latency: 10000µs
		ProxyEndTime:      now.Add(20 * time.Millisecond), // Total: 20000µs
		Success:           true,
	}

	history.AddRecord(record)

	// Verify the calculated metrics
	records := history.GetRecords()
	assert.Len(t, records, 1)

	calculatedRecord := records[0]
	assert.Equal(t, int64(20000), calculatedRecord.TotalDurationUs)   // 20ms = 20000µs total
	assert.Equal(t, int64(10000), calculatedRecord.UpstreamLatencyUs) // 10ms = 10000µs upstream
	assert.Equal(t, int64(10000), calculatedRecord.ProxyOverheadUs)   // 10ms = 10000µs proxy overhead

	// Verify stats calculation
	stats := history.GetStats()
	assert.Equal(t, int64(20000), stats["avg_duration_us"])
	assert.Equal(t, int64(10000), stats["avg_upstream_latency_us"])
	assert.Equal(t, int64(10000), stats["avg_proxy_overhead_us"])
}
