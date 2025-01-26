templ:
	templ generate --watch --cmd="go run . --serve" --proxy="http://localhost:8080"
tailwind:
	bunx @tailwindcss/cli -i ./public/css/tw-input.css -o ./public/css/tw-output.css --watch
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate
dev: 
	make -j3 sqlc tailwind templ 
