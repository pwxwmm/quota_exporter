### **新增 `README.md`**
```md
# Quota Exporter

## **Overview**
`quota_exporter` is a Prometheus exporter that collects storage quota usage metrics from a configured API and exposes them via HTTP.

## **Installation**
### **1. Clone the repository**
```bash
git clone ssh://git@code.iflytek.com:30004/mmwei3/quota_exporter.git
cd quota_exporter
```

### **2. Build the binary**
```bash
go build -o quota_exporter main.go
```

### **3. Run the exporter**
```bash
./quota_exporter
```

## **Configuration**
Edit `config/config.yaml`:
```yaml
quota_api_url: "http://example.com/api/quota"
```

## **Usage**
### **Prometheus Scraping**
Add the following job to your Prometheus configuration:
```yaml
scrape_configs:
  - job_name: "quota_exporter"
    static_configs:
      - targets: ["localhost:9533"]
```

### **Verify Metrics**
```bash
curl http://localhost:9533/metrics
```

## **Error Handling**
If the API fails, the metrics will be set to `0`.

## **Author**
- **Name:** mmwei3
- **Email:** mmwei3@iflytek.com, 1300042631@qq.com
- **Date:** 2025-03-20
```

---
