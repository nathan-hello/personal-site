templ:
	templ generate

air:
	go run github.com/cosmtrek/air@v1.52.0

tailwind:
	bunx tailwindcss -i ./public/css/tw-input.css -o ./public/css/tw-output.css --watch

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate
serve:
	go run main.go --serve
dev: 
	make -j4 templ air tailwind serve
