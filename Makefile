js/install:
	bun install
js/tailwind:
	bun run tailwindcss -i ./public/css/tw-input.css -o ./public/css/tw-output.css
js/tailwind/watch:
	bunx tailwindcss -i ./public/css/tw-input.css -o ./public/css/tw-output.css --watch
js/prod:
	make js/install js/tailwind

build/templ:
	templ generate
build/sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate
build/go:
	go build -o personal-site .

build:
	make build/sqlc build/css build/templ build/go

run/go:
	go run . --dev & echo $! > personal-site.pid

start:
	./cicd.sh

watch/templ:
	templ generate --watch --cmd="go run . --dev" --proxy="http://localhost:3000" --open-browser=false
dev: 
	make -j3 build/sqlc build/css build/templ asdf/go

prod:
	make build/sqlc build/templ build/go start
