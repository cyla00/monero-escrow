dev:
	@templ generate
	@go run ./cmd/main.go

pod up:
	@podman compose up

pod down:
	@podman compose down

templ:
	@templ generate