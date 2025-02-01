build/templ:
	templ generate
build/css:
	bunx @tailwindcss/cli -i ./public/css/tw-input.css -o ./public/css/tw-output.css
build/sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate


watch/templ:
	templ generate --watch --cmd="go run . --serve" --proxy="http://localhost:3000" --open-browser=false
watch/css:
	bunx tailwindcss -i ./public/css/tw-input.css -o ./public/css/tw-output.css --watch
dev: 
	make -j3 build/sqlc watch/css watch/templ 
