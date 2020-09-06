run-dev:
	docker-compose up -d
	sleep 10
	HTTP_PORT=8080 DB_DSN="user:password@tcp(127.0.0.1:3306)/db" go run ./cmd/api

.PHONY: run-dev
