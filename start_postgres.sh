#!/bin/bash

docker compose down -v
docker compose up -d --build

echo "Ждём пока Postgres поднимется..."
until docker exec my_postgres pg_isready -U postgres; do
    sleep 1
done

echo "Готово. Всё поднялось!"
