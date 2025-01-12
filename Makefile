dev/templ:
	templ generate --watch

dev/air:
	go run github.com/cosmtrek/air@v1.52.0

dev/tailwind:
	bunx tailwindcss -i ./public/css/tw-input.css -o ./public/css/tw-output.css --watch

dev/sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate

dev: 
	make -j3 dev/templ dev/air dev/tailwind
