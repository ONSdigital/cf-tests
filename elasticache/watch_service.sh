#!/usr/bin/env bash
set -eu
: ${ELASTICACHE_SERVICE_NAME:=test-elasticache}

get_status() {
    cf service "$ELASTICACHE_SERVICE_NAME" | grep Status:
}

while :; do
    status=$(get_status)
    [ "$status" = 'Status: create in progress' ] || break
    sleep 1
done

[ "$status" = 'Status: create succeeded' ] && exit 0

cf service "$ELASTICACHE_SERVICE_NAME"
exit 1