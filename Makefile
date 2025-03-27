build:
	go build -o quota_exporter ./cmd/main.go

run:
	./quota_exporter

docker-build:
	docker build -t quota_exporter .

docker-run:
	docker run -p 8866:8866 quota_exporter