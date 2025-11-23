# Quick Start Guide

Быстрый запуск и тестирование сервиса в 3 шага.

## Шаг 1: Запустить БД и применить миграции

```bash
# Запустить PostgreSQL
docker-compose up -d

# Дождаться готовности БД (5-10 секунд)
sleep 10

# Применить миграции
make migrate-up
```

## Шаг 2: Запустить сервис

В одном терминале:

```bash
make run
```

Вы должны увидеть:
```
INFO: Starting HTTP server on port 8080
INFO: Server is ready to handle requests at http://localhost:8080
```

## Шаг 3: Протестировать API

В другом терминале:

```bash
./test-api.sh
```

Скрипт проверит все эндпоинты и покажет результаты с ✓ или ✗.

## Быстрая проверка вручную

```bash
# Health check
curl http://localhost:8080/health

# Создать команду
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{"team_name":"test","members":[{"user_id":"u1","username":"Alice","is_active":true}]}'

# Получить команду
curl "http://localhost:8080/team/get?team_name=test" | jq .
```

## Остановка

```bash
# Остановить сервис: Ctrl+C (graceful shutdown)

# Остановить БД
docker-compose down
```

## Полная документация

См. [TESTING.md](TESTING.md) для детального тестирования.
