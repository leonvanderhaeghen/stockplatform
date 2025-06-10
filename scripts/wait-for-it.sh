#!/bin/sh
# wait-for-it.sh
# This script waits for a service to be available before continuing.

set -e

host="$1"
shift
port="$1"
shift
timeout="${WAITFORIT_TIMEOUT:-15}"
shift
cmd="$@"

if ! [ "$host" -a "$port" ]; then
    echo "Usage: wait-for-it.sh host port [-- command args]"
    exit 1
fi

echo "Waiting for $host:$port to be available..."

start_ts=$(date +%s)
while :
do
    (echo > /dev/tcp/$host/$port) >/dev/null 2>&1
    result=$?
    if [ $result -eq 0 ]; then
        end_ts=$(date +%s)
        echo "$host:$port is available after $((end_ts - start_ts)) seconds"
        break
    fi
    echo "Waiting for $host:$port..."
    sleep 1
    current_ts=$(date +%s)
    if [ $((current_ts - start_ts)) -ge $timeout ]; then
        echo "Timed out waiting for $host:$port"
        exit 1
    fi
done

exec $cmd
