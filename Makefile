# Makefile
export PATH := $(PATH):$(HOME)/go/bin

include .env
export

APP_NAME=go-task-service
MAIN_FILE=cmd/main.go
MIGRATIONS_DIR=migrations
SEARCH_DIRS=./cmd,./internal
DOCS_DIR=internal/docs
LOG_PREFIX=[swagger]
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
SEED_DIR=seed

.PHONY: help

run: ## Запуск сервера
	go run $(MAIN_FILE)

seed-tasks: ## Генерация задач в базу
	@echo "[seed] 🚜 Заполняем таблицу 'task' 1000 записями..."
	@psql $(DB_URL) -f $(SEED_DIR)/tasks.sql
	@echo "[seed] ✅ Генерация завершена."

migrate-reset:
	@echo "[migrate] 🔍 Проверяем текущую версию..."
	@VERSION_RAW=$$(migrate -source file://$(MIGRATIONS_DIR) -database "$(DB_URL)" version 2>&1); \
	echo "[migrate] 🧾 Raw output: $$VERSION_RAW"; \
	if echo "$$VERSION_RAW" | grep -qi dirty; then \
		VERSION=$$(echo "$$VERSION_RAW" | grep -oE '[0-9]+' | head -1); \
		echo "[migrate] ⚠️ База грязная на версии $$VERSION, сбрасываем..."; \
		migrate -source file://$(MIGRATIONS_DIR) -database "$(DB_URL)" force $$VERSION; \
		echo "[migrate] ✅ Dirty-флаг снят, можно продолжать миграции."; \
	else \
		echo "[migrate] ✅ База чистая, сброс не требуется."; \
	fi


swagger: check-swagger-deps check-annotations
	@echo "$(LOG_PREFIX) 🔥 Генерация Swagger документации..."
	@swag init -g $(MAIN_FILE) -d .
	@echo "$(LOG_PREFIX) ✅ Swagger готов: http://localhost:8080/swagger/index.html"

clean-swagger:
	@echo "$(LOG_PREFIX) 🧹 Удаляем старые Swagger-файлы..."
	@rm -rf $(DOCS_DIR)
	@echo "$(LOG_PREFIX) 🧼 $(DOCS_DIR) очищено."

check-swagger-deps:
	@test -f go.mod || (echo "$(LOG_PREFIX) ❌ Файл go.mod не найден!" && exit 1)
	@test -f $(MAIN_FILE) || (echo "$(LOG_PREFIX) ❌ Не найден $(MAIN_FILE)!" && exit 1)
	@command -v swag >/dev/null || (echo "$(LOG_PREFIX) ❌ swag не установлен. go install github.com/swaggo/swag/cmd/swag@latest" && exit 1)

check-annotations:
	@grep -q "@title" $(MAIN_FILE) || (echo "$(LOG_PREFIX) ⚠️ В $(MAIN_FILE) нет аннотации @title" && exit 1)
	@grep -q "@version" $(MAIN_FILE) || (echo "$(LOG_PREFIX) ⚠️ В $(MAIN_FILE) нет аннотации @version" && exit 1)
	@grep -q "@description" $(MAIN_FILE) || (echo "$(LOG_PREFIX) ⚠️ В $(MAIN_FILE) нет аннотации @description" && exit 1)

swag: ## Генерация swagger
	swag init -g cmd/main.go -d .

migrate-up: ## Запуск миграции
	migrate -source file://$(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -source file://$(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-force:
	migrate -source file://$(MIGRATIONS_DIR) -database "$(DB_URL)" force 1

migrate-version:
	migrate -source file://$(MIGRATIONS_DIR) -database "$(DB_URL)" version

build:
	go build -o $(APP_NAME) $(MAIN_FILE)

clean:
	rm -f $(APP_NAME)
