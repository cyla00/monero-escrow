dev:
	templ generate && go run ./cmd/main.go

db-up:
	podman compose up

db-down:
	podman compose down

templ:
	templ generate