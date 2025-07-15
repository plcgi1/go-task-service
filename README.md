# Задание

Необходимо разработать полноценный микросервис на Go работающий в высоконагруженной системе.

Бизнес-логика:
В базе данных определена таблица task. 
Каждая запись содержит поле status со значениями: NEW, PROCESSING, PROCESSED. 
Изначальное значение статуса - NEW, финальное - PROCESSED.
Микросервис должен предоставлять эндпоинт, 
при обращении к которому будет браться в обработку заданное количество записей из таблицы task 
и производится обработка согласно описанию ниже.

Параметры эндпоинта:
1. Количество записей, которые нужно взять в работу из базы;
2. Диапазон времени (в миллисекундах), 
   в течение которого будет обрабатываться одна запись (случайное значение в заданном диапазоне). 
   Если параметры не заданы, обработка происходит без задержки;
3. Вероятность успешного получения финального статуса PROCESSED.

## Переменные окружения 

.env файл в репо сейчас - для быстрого запуска - для прода в репо его быть не должно

```
DB_HOST=localhost # для локалки
# DB_HOST=db # для докера
DB_PORT=5434 # для локалки
# DB_PORT=5432 # для докера
DB_USER=taskdb_user
DB_PASSWORD=secret
DB_NAME=taskdb
APP_PORT=8080 # порт микросервиса
WORKERS=5 # количество обработчиков
SSL_MODE=disable
COUNT_OF_TRYINGS=N # допусчтимое количество повторных запусков задачи
```

## Локальный запуск

```
// psql
CREATE USER taskdb_user WITH PASSWORD 'secret';
ALTER USER taskdb_user WITH SUPERUSER;
ALTER ROLE taskdb_user CREATEROLE CREATEDB;

CREATE DATABASE taskdb;
GRANT ALL PRIVILEGES ON DATABASE taskdb to taskdb_user;

// shell
make migrate-up # один раз
make seed-tasks # один раз - закинет в базу 1000 задач для тестов

make run # старт сервера
```

### Makefile команды

смотри Makefile 

## Docker запуск

1. docker compose build
2. docker compose up

http://localhost:8080/swagger/index.html - запуск свагера - через который делаем запросы

