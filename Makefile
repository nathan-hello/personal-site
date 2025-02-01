build/templ:
	templ generate
build/css:
	bunx @tailwindcss/cli -i ./public/css/tw-input.css -o ./public/css/tw-output.css
build/sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate


watch/templ:
	templ generate --watch --cmd="go run . --serve" --proxy="http://localhost:8080"
watch/tailwind:
	bunx @tailwindcss/cli -i ./public/css/tw-input.css -o ./public/css/tw-output.css --watch
dev: 
	make -j3 build/sqlc dev/tailwind dev/templ 
