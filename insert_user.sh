#!/bin/bash

echo "⏳ Добавляю тестового пользователя..."

docker exec -i my_postgres psql -U postgres -d students <<EOF
INSERT INTO users (name_, surname, email)
VALUES ('Sagir', 'Jusupov', 'alexo98@yandex.ru')
ON CONFLICT DO NOTHING;
EOF

echo "Пользователь добавлен (если ещё не был)"
