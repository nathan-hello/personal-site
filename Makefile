setup:
	bun install

build/templ:
	templ generate
build/css:
	bun run tailwindcss -i ./public/css/tw-input.css -o ./public/css/tw-output.css
build/sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate
build/go:
	go build -o personal-site .

build:
	make build/sqlc build/css build/templ build/go
