start:
	fuser -k 8000/tcp 2>/dev/null ; go run ./cmd/web

test:
	go test -v -cover ./internal/server -count=1

checkServer:
	curl http://localhost:8000/api/check_health
