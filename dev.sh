#!/bin/bash

set -e

APP_NAME="personal-site"
PORT=3000
TEMPL_DIRS="components layouts"
TAILWIND_INPUT="./public/css/tw-input.css"
TAILWIND_OUTPUT="./public/css/tw-output.css"
GO_FILES="."

cleanup() {
    echo "Cleaning up..."
    kill $TEMPL_PID $TAILWIND_PID $SERVER_PID 2>/dev/null || true
    kill $(lsof -t -i:$PORT) 2>/dev/null || true
}

trap cleanup EXIT INT TERM

echo "Starting Tailwind CSS in watch mode..."
bunx tailwindcss -i $TAILWIND_INPUT -o $TAILWIND_OUTPUT --watch &
TAILWIND_PID=$!

echo "Starting templ in watch mode..."
templ generate --watch &
TEMPL_PID=$!

start_server() {
    echo "Starting Go server..."
    go run . --dev &
    SERVER_PID=$!
}

start_server

echo "Watching for Go file changes..."
while true; do
    sleep 1
    
    NEWEST=$(find $GO_FILES -name "*.go" -not -name "*_templ.go" -not -name "*_gen.go" -newer /tmp/.dev_server_marker 2>/dev/null | head -1)
    
    if [ -n "$NEWEST" ]; then
        echo "Detected change in $NEWEST, restarting server..."
        kill $SERVER_PID 2>/dev/null || true
        wait $SERVER_PID 2>/dev/null || true
        
        if lsof -i:$PORT >/dev/null 2>&1; then
            kill $(lsof -t -i:$PORT) 2>/dev/null || true
            sleep 1
        fi
        
        touch /tmp/.dev_server_marker
        start_server
    fi
done
