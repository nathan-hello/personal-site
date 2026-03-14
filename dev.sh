#!/usr/bin/env bash

set -uo pipefail

APP_PORT="${APP_PORT:-3000}"
TEMPL_PROXY_PORT="${TEMPL_PROXY_PORT:-7331}"

SERVER_PID=""
TEMPL_PID=""
TAILWIND_PID=""

log() {
    printf '==> %s\n' "$*"
}

stop_pid() {
    local pid="${1:-}"

    if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
        kill "$pid" 2>/dev/null || true
        wait "$pid" 2>/dev/null || true
    fi
}

stop_pid_tree() {
    local pid="${1:-}"
    local child

    if [ -z "$pid" ] || ! kill -0 "$pid" 2>/dev/null; then
        return
    fi

    while IFS= read -r child; do
        [ -n "$child" ] && stop_pid_tree "$child"
    done < <(pgrep -P "$pid" 2>/dev/null || true)

    stop_pid "$pid"
}

kill_port() {
    local pids

    pids="$(lsof -ti tcp:"$APP_PORT" 2>/dev/null || true)"
    if [ -n "$pids" ]; then
        kill $pids 2>/dev/null || true
        sleep 1
    fi
}

stop_server() {
    stop_pid_tree "$SERVER_PID"
    SERVER_PID=""
    kill_port
}

start_server() {
    make dev/server &
    SERVER_PID=$!
}

notify_browser() {
    make dev/reload >/dev/null 2>&1 || true
}

run_make() {
    if make "$@"; then
        return 0
    fi

    log "build failed; keeping current server running"
    return 1
}

restart_server() {
    local reason="$1"
    shift

    log "$reason"
    if [ "$#" -gt 0 ] && ! run_make "$@"; then
        return
    fi

    stop_server
    start_server
    notify_browser
}

cleanup() {
    trap - EXIT INT TERM
    log "Stopping dev server"
    stop_pid_tree "$TAILWIND_PID"
    stop_pid_tree "$TEMPL_PID"
    stop_server
}

trap cleanup EXIT INT TERM

log "Running initial build"
if ! run_make dev/bootstrap; then
    exit 1
fi

log "Starting app server on :$APP_PORT"
kill_port
start_server

log "Starting templ proxy on :$TEMPL_PROXY_PORT"
make dev/templ &
TEMPL_PID=$!

log "Starting Tailwind watcher"
make dev/tailwind &
TAILWIND_PID=$!

sleep 2

log "Open http://127.0.0.1:$TEMPL_PROXY_PORT for live reload"
log "Watching Go, templ, pages, public content, and Tailwind output"

while IFS= read -r changed_path; do
    case "$changed_path" in
        ./dist/*|./.git/*|./node_modules/*)
            continue
            ;;
        *_templ.go|*_templ.txt|*_gen.go|*.sql.go|./personal-site)
            continue
            ;;
        ./public/css/tw-input.css)
            continue
            ;;
        ./public/css/tw-output.css)
            restart_server "tailwind rebuilt: ${changed_path#./}"
            ;;
        *.templ)
            restart_server "templ changed: ${changed_path#./}" dev/build-templ
            ;;
        *.go)
            restart_server "go changed: ${changed_path#./}" dev/build-go
            ;;
        *.sql|./sqlc.yml)
            restart_server "sql changed: ${changed_path#./}" dev/build-sql
            ;;
        *.html|*.mdx)
            restart_server "content changed: ${changed_path#./}"
            ;;
        ./pages/*|./public/*)
            restart_server "public asset changed: ${changed_path#./}"
            ;;
    esac
done < <(
    inotifywait -q -m -r \
        --format '%w%f' \
        --event close_write,create,move,delete \
        --exclude '(^|/)(dist|node_modules|\.git)(/|$)|(_templ\.go$|_templ\.txt$|_gen\.go$|\.sql\.go$)|(^|/)personal-site$' \
        .
)
