#!/bin/sh
set -e

host="$1"
port="${host##*:}"
host="${host%:*}"
shift
cmd="$@"

echo "Waiting for $host:$port..."

until nc -z "$host" "$port"; do
  echo "Waiting for $host:$port..."
  sleep 1
done

echo "$host:$port is available"
exec $cmd