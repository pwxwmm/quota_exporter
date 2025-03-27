package main

import (
	"context"
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
	Author  = "mmwei3"
	Email   = "mmwei3@iflytek.com, 1300042631@qq.com"
	Date    = "2025-03-25"
)

func main() {
	// Print author information
	fmt.Println("=====================================")
	fmt.Println("           Quota Exporter            ")
	fmt.Println("=====================================")
	fmt.Printf("Author: %s\nEmail: %s\nDate: %s\n", Author, Email, Date)

	// Load configuration
	cfg := config.LoadConfig("/root/config.yaml")
	if cfg.QuotaExporter.APIURL == "" {
		log.Fatal("Quota API URL is missing in config.yaml. Exiting...")
	}

	// Initialize Prometheus metrics
	collector.InitMetrics()

	// Start periodic metric updates
	go updateMetrics(cfg)

	// Start HTTP server for Prometheus scraping
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Quota Exporter is running on http://0.0.0.0:9533/metrics")
	log.Fatal(http.ListenAndServe(":9533", nil))
}

// Periodically update metrics
func updateMetrics(cfg config.Config) {
	for {
		var wg sync.WaitGroup
		// Use RequestTimeout from config file
		ctx, cancel := context.WithTimeout(context.Background(), cfg.QuotaExporter.RequestTimeout)
		defer cancel()

		wg.Add(1)
		go func() {
			defer wg.Done()
			// Pass the entire QuotaExporterConfig struct to FetchQuotaData
			collector.FetchQuotaData(ctx, cfg.QuotaExporter)
		}()

		wg.Wait()
		// Sleep for scrape_interval
		time.Sleep(time.Duration(cfg.QuotaExporter.ScrapeInterval) * time.Second)
	}
}
