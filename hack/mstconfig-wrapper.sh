#!/usr/bin/env bash

REAL=/usr/bin/mstconfig.real
TRIES=20

for i in $(seq 1 $TRIES); do
    out=$(mktemp)
    err=$(mktemp)

    "$REAL" "$@" >"$out" 2>"$err"
    rc=$?

    if [ $rc -ne 3 ]; then
        cat "$out"
        cat "$err" >&2
        rm -f "$out" "$err"
        exit $rc
    fi

    rm -f "$out" "$err"

    sleep 0.2
done

exit 3
