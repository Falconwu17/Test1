#!/usr/bin/env bash
URL="http://localhost:8080/entries"
TOTAL=10000

echo "Отправляем $TOTAL entry на $URL ..."

for ((i=1; i<=TOTAL; i++)); do
    val=$(echo "scale=5; 30 + 30 * s($i/100)" | bc -l | tr -d '\n')

    json=$(jq -n \
        --arg rid "1" \
        --arg v "$val" \
        '{record_id: ($rid | tonumber), data: {value: ($v | tonumber)}}')

    curl -s -X POST "$URL" \
         -H "Content-Type: application/json" \
         -d "$json" >/dev/null
    sleep 0.05
    (( i % 100 == 0 )) && echo "Отправлено $i..."
done

echo "Готово!"
