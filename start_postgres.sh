#!/bin/bash

docker compose down -v
docker compose up -d

docker logs my_postgres

echo ждем пока все наладиться
until docker exec my_postgres pg_isready -U postgres; do
    sleep 1
done

echo "Postgres поднят, ждём ещё 3 секунды , пока применится init.sql"
sleep 3

echo "Запуск Go"
go run ./cmd/app
