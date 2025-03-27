// Package collector implements quota metrics collection for Prometheus.
//
// Author:  mmwei3
// Email:   mmwei3@iflytek.com, 1300042631@qq.com
// Date:    2025-03-25

package collector

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"quota_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// Mock quota API server
func mockQuotaServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	}))
}

// Test FetchQuotaData function
func TestFetchQuotaData(t *testing.T) {
	mockResponse := `{"status": true, "result": {"quotas": [{"path": "/group1/user1", "size_sum": 1000, "size_used": 500}]}}`
	server := mockQuotaServer(mockResponse, http.StatusOK)
	defer server.Close()

	// Define storage configuration based on the new structure
	quotaConfig := config.QuotaExporterConfig{
		APIURL:        server.URL, // Now using api_url from config
		ScrapeInterval: 600,
	}

	// Initialize metrics
	quotaMetrics := make(map[string]*prometheus.GaugeVec)
	quotaMetrics["available"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "quota_size_available",
			Help: "Available quota size",
		},
		[]string{"group", "user"},
	)
	quotaMetrics["used"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "quota_size_used",
			Help: "Used quota size",
		},
		[]string{"group", "user"},
	)
	quotaMetrics["free"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "quota_size_free",
			Help: "Free quota size",
		},
		[]string{"group", "user"},
	)

	// Register metrics
	prometheus.MustRegister(quotaMetrics["available"])
	prometheus.MustRegister(quotaMetrics["used"])
	prometheus.MustRegister(quotaMetrics["free"])

	// Simulate FetchQuotaData function, assuming it uses quotaConfig to fetch the data
	FetchQuotaData(quotaConfig, quotaMetrics)

	// Verify metrics
	available := testutil.ToFloat64(quotaMetrics["available"].WithLabelValues("group1", "user1"))
	used := testutil.ToFloat64(quotaMetrics["used"].WithLabelValues("group1", "user1"))
	free := testutil.ToFloat64(quotaMetrics["free"].WithLabelValues("group1", "user1"))

	if available != 1000 {
		t.Errorf("Expected available quota 1000, got %f", available)
	}
	if used != 500 {
		t.Errorf("Expected used quota 500, got %f", used)
	}
	if free != 500 {
		t.Errorf("Expected free quota 500, got %f", free)
	}
}
