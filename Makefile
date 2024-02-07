dev:
	@templ generate
	@go run ./cmd/main.go

up:
	@podman compose up

down:
	@podman compose down

templ:
	@templ generate