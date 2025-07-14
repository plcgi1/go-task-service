#!/bin/bash
set -e

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${SSL_MODE}"

echo "[docker] 📤 Запускаем миграции..."
migrate -source file://migrations -database "$DATABASE_URL" up

echo "[docker] 🚜 Seed задач..."
psql "$DATABASE_URL" -f seed/tasks.sql

echo "[docker] 🏁 Запуск приложения..."
exec ./app