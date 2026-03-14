.PHONY: js/install js/tailwind js/tailwind/watch build/css build/templ build/sqlc build/go build serve dev dev/bootstrap dev/server dev/build-go dev/build-templ dev/build-sql dev/templ dev/tailwind dev/reload

APP_BIN := personal-site
APP_PORT ?= 3000
TEMPL_PROXY_PORT ?= 7331
TAILWIND_INPUT := ./public/css/tw-input.css
TAILWIND_OUTPUT := ./public/css/tw-output.css

js/install:
	bun install

js/tailwind:
	bunx tailwindcss -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT)

js/tailwind/watch:
	bunx tailwindcss -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT) --watch

build/css: js/install js/tailwind

build/templ:
	templ generate

build/sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate

build/go:
	go build -o $(APP_BIN) .

build:
	$(MAKE) build/sqlc build/css build/templ build/go

# don't forget to allow nginx to read from dist/public
# doas chown -R morrow:nginx ./dist
serve: build
	./$(APP_BIN) --build --serve

dev: dev/bootstrap
	./dev.sh

dev/bootstrap: build

dev/server:
	./$(APP_BIN) --dev

dev/build-go:
	$(MAKE) build/go

dev/build-templ:
	$(MAKE) build/templ build/go

dev/build-sql:
	$(MAKE) build/sqlc build/go

dev/templ:
	templ generate --watch --proxy="http://127.0.0.1:$(APP_PORT)" --proxyport="$(TEMPL_PROXY_PORT)" --open-browser=false

dev/tailwind:
	bunx tailwindcss -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT) --watch

dev/reload:
	templ generate --notify-proxy --proxyport="$(TEMPL_PROXY_PORT)"
