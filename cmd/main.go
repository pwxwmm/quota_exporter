package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"quota_exporter/collector"
	"quota_exporter/config"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ============== Author Information ==============
const (
	Author = "mmwei3"
	Email  = "mmwei3@iflytek.com, 1300042631@qq.com"
	Date   = "2025-03-25"
)

func main() {
	fmt.Println("=====================================")
	fmt.Println("           Quota Exporter            ")
	fmt.Println("=====================================")
	fmt.Printf("Author: %s\nEmail: %s\nDate: %s\n", Author, Email, Date)

	cfg := config.LoadConfig("/root/config.yaml")
	if cfg.QuotaExporter.APIURL == "" {
		log.Fatal("Quota API URL is missing in config.yaml. Exiting...")
	}

	collector.InitMetrics()
	go updateMetrics(cfg)

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		if collector.GetHealthStatus() != 1 {
			status = "failed"
		}
		resp := map[string]interface{}{
			"status":       status,
			"last_updated": collector.GetLastUpdated().Format("2006-01-02 15:04:05"),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("Quota Exporter is running on http://0.0.0.0:9533/metrics")
	log.Fatal(http.ListenAndServe(":9533", nil))
}

func updateMetrics(cfg config.Config) {
	for {
		var wg sync.WaitGroup
		ctx, cancel := context.WithTimeout(context.Background(), cfg.QuotaExporter.RequestTimeout)
		defer cancel()

		wg.Add(1)
		go func() {
			defer wg.Done()
			collector.FetchQuotaData(ctx, cfg.QuotaExporter)
		}()

		wg.Wait()
		time.Sleep(time.Duration(cfg.QuotaExporter.ScrapeInterval) * time.Second)
	}
}
