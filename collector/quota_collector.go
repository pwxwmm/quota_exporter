package collector

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"quota_exporter/config"
	"strconv"
	"strings"
	"sync"
	"time"
	"math"

	"github.com/prometheus/client_golang/prometheus"
)

type QuotaAPIResponse struct {
	Code string         `json:"code"`
	Data []QuotaMetrics `json:"data"`
}

type QuotaMetrics struct {
	SizeUsed  int64  `json:"size_used"`
	CluName   string `json:"clu_name"`
	UseRate   string `json:"use_rate"`
	GroupName string `json:"group_name"`
	Date      string `json:"date"`
	SizeSum   int64  `json:"size_sum"`
}

var (
	quotaMetrics         = map[string]*prometheus.GaugeVec{}
	lastUpdatedGauge     prometheus.Gauge
	exporterHealthGauge  prometheus.Gauge
	lastUpdatedTime      time.Time
	lastUpdatedTimeMutex sync.RWMutex
	healthStatus         = 1
)

func InitMetrics() {
	quotaMetrics["quota_group_use_ratio"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "quota_group_use_ratio",
			Help: "Quota group use rate percentage",
		},
		[]string{"storage", "group_name"},
	)

	quotaMetrics["cluster_use_rate"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cluster_use_rate",
			Help: "Cluster use rate percentage",
		},
		[]string{"storage"},
	)

	lastUpdatedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "quota_last_updated_timestamp",
			Help: "Last successful quota update timestamp (Unix seconds)",
		},
	)

	exporterHealthGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "exporter_health_status",
			Help: "Exporter health status: 1 = ok, 0 = failed",
		},
	)

	for _, metric := range quotaMetrics {
		prometheus.MustRegister(metric)
	}
	prometheus.MustRegister(lastUpdatedGauge)
	prometheus.MustRegister(exporterHealthGauge)
}

func FetchQuotaData(ctx context.Context, cfg config.QuotaExporterConfig) {
	startTime := time.Now()
	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, "POST", cfg.APIURL, nil)
	if err != nil {
		log.Printf("[ERROR] %s - Failed to create request: %v", formatTimestamp(startTime), err)
		setExporterHealth(0)
		setQuotaToZero()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] %s - Failed to fetch quota data: %v", formatTimestamp(startTime), err)
		setExporterHealth(0)
		setQuotaToZero()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] %s - Quota API returned status: %d", formatTimestamp(startTime), resp.StatusCode)
		setExporterHealth(0)
		setQuotaToZero()
		return
	}

	var result QuotaAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[ERROR] %s - Failed to parse JSON: %v", formatTimestamp(startTime), err)
		setExporterHealth(0)
		setQuotaToZero()
		return
	}

	if result.Code == "0000" {
		log.Printf("[INFO] %s - Successfully fetched quota data. Entries: %d", formatTimestamp(startTime), len(result.Data))
		processQuotaData(result.Data)
		updateLastUpdated(time.Now())
		setExporterHealth(1)
	} else {
		log.Printf("[ERROR] %s - API returned unexpected code: %s", formatTimestamp(startTime), result.Code)
		setExporterHealth(0)
		setQuotaToZero()
	}
}

func processQuotaData(quotas []QuotaMetrics) {
	for _, quota := range quotas {
		if quota.GroupName == "" {
			quota.GroupName = "unknown"
		}

		useRateValue := parsePercentage(quota.UseRate)
		labelValues := []string{quota.CluName, quota.GroupName}
		clusterLabel := []string{quota.CluName}

		ratio := float64(quota.SizeUsed) / float64(quota.SizeSum) * 100
		ratio = math.Round(ratio*100) / 100 // 保留两位小数

		quotaMetrics["quota_group_use_ratio"].WithLabelValues(labelValues...).Set(ratio)
		quotaMetrics["cluster_use_rate"].WithLabelValues(clusterLabel...).Set(useRateValue)
	}
}

func setQuotaToZero() {
	defaultLabels := []string{"unknown", "unknown"}
	clusterLabels := []string{"unknown"}

	quotaMetrics["quota_group_use_ratio"].WithLabelValues(defaultLabels...).Set(0)
	quotaMetrics["cluster_use_rate"].WithLabelValues(clusterLabels...).Set(0)
	updateLastUpdated(time.Unix(0, 0))
}

func parsePercentage(value string) float64 {
	cleanValue := strings.TrimSuffix(value, "%")
	parsedValue, err := strconv.ParseFloat(cleanValue, 64)
	if err != nil {
		log.Printf("[ERROR] Failed to parse use_rate '%s': %v", value, err)
		return 0
	}
	return parsedValue
}

func formatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Prometheus + /health 用
func updateLastUpdated(t time.Time) {
	lastUpdatedTimeMutex.Lock()
	defer lastUpdatedTimeMutex.Unlock()
	lastUpdatedTime = t
	lastUpdatedGauge.Set(float64(t.Unix()))
}

func GetLastUpdated() time.Time {
	lastUpdatedTimeMutex.RLock()
	defer lastUpdatedTimeMutex.RUnlock()
	return lastUpdatedTime
}

func setExporterHealth(status int) {
	exporterHealthGauge.Set(float64(status))
	healthStatus = status
}

func GetHealthStatus() int {
	return healthStatus
}
