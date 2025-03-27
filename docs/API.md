### **API.md**

# **API Documentation for quota_exporter**

## **1. Overview**
`quota_exporter` exposes storage quota data via an HTTP endpoint for Prometheus scraping. It uses a configurable API to pull data from different storage backends.

---

## **2. API Endpoint**

### **GET /metrics**

**Description**:  
Fetches the storage quota metrics in a format suitable for Prometheus scraping.

**URL**:  
`http://<hostname>:9533/metrics`

**Response Format**:  
The response contains metrics in the following Prometheus format:

```plaintext
# HELP quota_size_available Quota size for storage
# TYPE quota_size_available gauge
quota_size_available{storage="train28", group_name: "bitbrain",date: "2025-01-20 17:27"} 1050000
quota_size_used{storage="train28", group_name: "bitbrain", date: "2025-01-20 17:27"} 500000
quota_size_free{storage="train28", group_name: "bitbrain", date: "2025-01-20 17:27"} 550000
cluster_use_rate{storage="train28", date: "2025-01-20 17:27"} 42,
```

**Prometheus Scrape Configuration Example**:  
Add the following to your Prometheus configuration:

```yaml
scrape_configs:
  - job_name: 'quota_exporter'
    static_configs:
      - targets: ['localhost:9533']
```

**Expected Metrics**:  
- `quota_size_available`: Total available quota.
- `quota_size_used`: Total used quota.
- `quota_size_free`: Total free quota.

--- 