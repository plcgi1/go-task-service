FROM golang:1.23

RUN apt-get update && apt-get install -y postgresql-client

WORKDIR /app

# Установка переменной окружения
ENV GOPROXY=direct

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Устанавливаем migrate
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Копируем остальной код
COPY . .

# Сборка приложения
RUN go build -o app ./cmd/main.go

COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh
RUN ls -la /app

ENTRYPOINT ["/app/start.sh"]
