#!/bin/bash
set -f

patterns=( '*.html' '*.go' '*.templ' )
cmd="make build/css build/templ run/go"

find_args=()
for pat in "${patterns[@]}"; do
    find_args+=( -name "$pat" -o )
done
unset 'find_args[-1]'

get_checksum() {
    files=$(find . -type f \( "${find_args[@]}" \) | sort)
    if [ -z "$files" ]; then
        echo ""
    else
        echo "$files" | xargs md5sum | md5sum
    fi
}

checksum=$(get_checksum)
if [ -z "$checksum" ]; then
    echo "No files found with pattern: $pattern"
    exit 1
fi

while true; do
    inotifywait -e modify,create,delete,move -r .
    new_checksum=$(get_checksum)
    if [ "$checksum" != "$new_checksum" ]; then
        if [ -f personal-site.pid ]; then
            pid=$(cat personal-site.pid)
            kill -9 $pid
            rm personal-site.pid
        fi
        $cmd
        checksum="$new_checksum"
    fi
done
