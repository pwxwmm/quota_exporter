// Package collector implements quota metrics collection for Prometheus.
//
// Author:  mmwei3
// Email:   mmwei3@iflytek.com, 1300042631@qq.com
// Date:    2025-03-25

package collector

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"quota_exporter/config"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)


// QuotaAPIResponse represents the API response structure
type QuotaAPIResponse struct {
	Code string         `json:"code"`
	Data []QuotaMetrics `json:"data"`
}

// QuotaMetrics represents individual quota metrics
type QuotaMetrics struct {
	SizeUsed  int64  `json:"size_used"`
	CluName   string `json:"clu_name"`
	UseRate   string `json:"use_rate"`
	GroupName string `json:"group_name"`
	Date      string `json:"date"`
	SizeSum   int64  `json:"size_sum"`
}

// Prometheus metrics map
var quotaMetrics = map[string]*prometheus.GaugeVec{}

// InitMetrics initializes Prometheus metrics
func InitMetrics() {
	quotaMetrics["quota_size_available"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "quota_size_available",
			Help: "Available quota size",
		},
		[]string{"storage", "group_name", "date"},
	)
	quotaMetrics["quota_size_used"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "quota_size_used",
			Help: "Used quota size",
		},
		[]string{"storage", "group_name", "date"},
	)
	quotaMetrics["quota_size_free"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "quota_size_free",
			Help: "Free quota size",
		},
		[]string{"storage", "group_name", "date"},
	)
	quotaMetrics["cluster_use_rate"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cluster_use_rate",
			Help: "Cluster use rate percentage",
		},
		[]string{"storage", "date"},
	)

	// Register metrics with Prometheus
	for _, metric := range quotaMetrics {
		prometheus.MustRegister(metric)
	}
}

// StartQuotaCollector runs a scheduled job to fetch quota data periodically
func StartQuotaCollector(ctx context.Context, cfg config.QuotaExporterConfig) {
	ticker := time.NewTicker(cfg.ScrapeInterval * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			FetchQuotaData(ctx, cfg)
		case <-ctx.Done():
			log.Println("[INFO] Stopping quota collector")
			return
		}
	}
}

// FetchQuotaData retrieves and processes quota data from the API
func FetchQuotaData(ctx context.Context, cfg config.QuotaExporterConfig) {
	startTime := time.Now()
	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, "POST", cfg.APIURL, nil)
	if err != nil {
		log.Printf("[ERROR] %s - Failed to create request: %v", formatTimestamp(startTime), err)
		setQuotaToZero()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] %s - Failed to fetch quota data: %v", formatTimestamp(startTime), err)
		setQuotaToZero()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] %s - Quota API returned status: %d", formatTimestamp(startTime), resp.StatusCode)
		setQuotaToZero()
		return
	}

	var result QuotaAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[ERROR] %s - Failed to parse JSON: %v", formatTimestamp(startTime), err)
		setQuotaToZero()
		return
	}

	if result.Code == "0000" {
		log.Printf("[INFO] %s - Successfully fetched quota data. Entries: %d", formatTimestamp(startTime), len(result.Data))
		processQuotaData(result.Data)
	} else {
		log.Printf("[ERROR] %s - API returned unexpected code: %s", formatTimestamp(startTime), result.Code)
		setQuotaToZero()
	}
}

// formatTimestamp returns a formatted timestamp string
func formatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// processQuotaData updates Prometheus metrics with fetched quota data
func processQuotaData(quotas []QuotaMetrics) {
	for _, quota := range quotas {
		// if group_name is null，set default unknow
		if quota.GroupName == "" {
			quota.GroupName = "unknown"
		}

		// trans use_rate delete "%"
		useRateValue := parsePercentage(quota.UseRate)

		// Prometheus label
		labelValues := []string{quota.CluName, quota.GroupName, quota.Date}
		clusterLabelValues := []string{quota.CluName, quota.Date}

		// set metrics value
		quotaMetrics["quota_size_available"].WithLabelValues(labelValues...).Set(float64(quota.SizeSum))
		quotaMetrics["quota_size_used"].WithLabelValues(labelValues...).Set(float64(quota.SizeUsed))
		quotaMetrics["quota_size_free"].WithLabelValues(labelValues...).Set(float64(quota.SizeSum - quota.SizeUsed))
		quotaMetrics["cluster_use_rate"].WithLabelValues(clusterLabelValues...).Set(useRateValue)
	}
}

// setQuotaToZero resets metrics to zero when API call fails
func setQuotaToZero() {
	defaultLabels := []string{"unknown", "unknown", "unknown"}
	clusterLabels := []string{"unknown", "unknown"}

	quotaMetrics["quota_size_available"].WithLabelValues(defaultLabels...).Set(0)
	quotaMetrics["quota_size_used"].WithLabelValues(defaultLabels...).Set(0)
	quotaMetrics["quota_size_free"].WithLabelValues(defaultLabels...).Set(0)
	quotaMetrics["cluster_use_rate"].WithLabelValues(clusterLabels...).Set(0)
}

// parsePercentage removes the '%' from a string and converts it to float64
func parsePercentage(value string) float64 {
	cleanValue := strings.TrimSuffix(value, "%")
	parsedValue, err := strconv.ParseFloat(cleanValue, 64)
	if err != nil {
		log.Printf("[ERROR] Failed to parse use_rate '%s': %v", value, err)
		return 0
	}
	return parsedValue
}
