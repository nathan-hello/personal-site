js/install:
	bun install
js/tailwind:
	bun run tailwindcss -i ./public/css/tw-input.css -o ./public/css/tw-output.css
js/tailwind/watch:
	bunx tailwindcss -i ./public/css/tw-input.css -o ./public/css/tw-output.css --watch

build/css:
	make js/install js/tailwind
build/templ:
	templ generate
build/sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate
build/go:
	go build -o personal-site .

build:
	make build/sqlc build/css build/templ build/go

# don't forget to allow nginx to read from dist/public
# doas chown -R morrow:nginx ./dist
serve:
	make build && ./personal-site --build --serve
