
---

# CHANGELOG
## [v1.2.2] - 2025-04-10

### Added
- Support for configurable `log_level` in `config.yaml` (`info` or `debug`).
- Debug logs for quota details will only be printed when `log_level: debug`.

### Changed
- Centralized log level control via `config.IsDebugEnabled()` utility function.


## [v1.2.1] - 2025-04-10

### Added
- Support for configurable `log_level` in `config.yaml` (`info` or `debug`).
- Debug logs for quota details will only be printed when `log_level: debug`.

### Changed
- Centralized log level control via `config.IsDebugEnabled()` utility function.

---

## [v1.2.0] - 2025-04-10

### Added
- `/health` endpoint returns JSON status with `status` and `last_updated`.
- Prometheus metric `exporter_health_status` to reflect exporter health (`1` for ok, `0` for failed).
- Prometheus metric `quota_last_updated_timestamp` to track the last successful scrape time.
- Prometheus metric `collector_health_status` for each collector's health status.

### Changed
- Refactored `/health` logic to use centralized status tracking.
- Precision of `quota_group_use_ratio` is now rounded to two decimal places.

---

## [v1.1.0] - 2025-04-01

### Added
- Prometheus metric `cluster_use_rate` for cluster-wide quota usage percentage.
- Dynamic label support: `storage`, `group_name`, and `date`.

### Changed
- Deprecated metrics `quota_size_available`, `quota_size_used`, and `quota_size_free` (commented out).
- Code cleanup and better logging format with timestamps.

---

## [v1.0.0] - 2025-03-25

### Initial Release
- Basic Prometheus exporter framework.
- Collects quota usage data from internal API.
- Exposes `quota_group_use_ratio` metric.

---