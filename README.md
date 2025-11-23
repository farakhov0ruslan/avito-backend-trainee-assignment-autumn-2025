
# PR Reviewer Assignment Service

Сервис для автоматического назначения ревьюеров на Pull Request’ы внутри команды.

---

## Описание

PR Reviewer Assignment Service — это HTTP API-сервис, который помогает автоматически назначать ревьюеров на pull request’ы. Он умеет:

- создавать и хранить команды и их участников;
- управлять активностью пользователей (участвуют ли в ротации ревью или нет);
- автоматически назначать до двух ревьюеров на каждый новый PR;
- перераспределять ревьюера (reassign), если нужно передать ревью другому человеку;
- выдавать список PR’ов, назначенных конкретному пользователю;
- отслеживать состояние PR (OPEN/MERGED);
- обеспечивать детерминированное и справедливое распределение ревью.

---

## Технологический стек

- **Go 1.24** — основной язык реализации;
- **PostgreSQL 15** — основная база данных;
- **gorilla/mux** — HTTP-роутер;
- **pgx/v5** — драйвер PostgreSQL и пул соединений;
- **goose** — миграции БД;
- **Docker & Docker Compose** — упаковка и запуск сервиса.

---

## Архитектура

Используется упрощённая Clean Architecture с разделением на слои:

```

Domain Models (internal/domain/models)
↓
Repository (internal/repository)
↓
Service (internal/service)
↓
Handler (internal/handler)
↓
HTTP Server

````

Основные уровни:

- **Domain Models** — сущности предметной области (`Team`, `User`, `PullRequest`, `PRReviewer` и т.п.).
- **Repository** — работа с базой данных (интерфейсы + реализация под PostgreSQL).
- **Service** — бизнес-логика (назначение ревьюеров, reassign, проверка состояний).
- **Handler** — HTTP-обработчики, преобразующие запрос/ответ в DTO.
- **HTTP Server** — настройка роутов, middleware, graceful shutdown.

---

## Требования

### Для запуска через Docker

- Docker
- Docker Compose
- Make (для удобного запуска команд, опционально)

### Для локальной разработки

- Go 1.24+
- PostgreSQL 15+
- Утилита **goose** для миграций
- Make (опционально)

---

## Запуск через Docker Compose

3. Поднять сервисы:

   ```bash
   docker-compose up --build
   ```

Приложение будет доступно по адресу:
`http://localhost:8080`

### Проверка health-эндпоинта

```bash
curl http://localhost:8080/health
# Ожидаемый ответ: HTTP 200 OK
```

---

## Локальный запуск (без Docker)

### Установка зависимостей

```bash
# Установить goose
make install-goose

# Установить дополнительные инструменты (линтеры и т.п.)
make install-tools

# Скачать Go-зависимости
go mod download
```

### Подготовка базы данных

1. Запустить PostgreSQL (локально или через Docker):

   ```bash
   make docker-up
   ```

2. Применить миграции:

   ```bash
   make migrate-up
   ```

3. Проверить статус миграций:

   ```bash
   make migrate-status
   ```

### Запуск приложения

```bash
# Через make
make run

# Либо напрямую
go run cmd/api/main.go
```

---

## API Endpoints

Полное описание API можно найти в файле `claude/openapi.yml`.

### Health Check

```http
GET /health
```

### Команды (teams)

**Создание команды:**

```http
POST /team/add
Content-Type: application/json

{
  "team_name": "backend-team",
  "members": [
    {"user_id": "user-1", "username": "Alice", "is_active": true},
    {"user_id": "user-2", "username": "Bob", "is_active": true}
  ]
}
```

**Получение команды:**

```http
GET /team/get?team_name=backend-team
```

---

### Пользователи (users)

**Изменение активности пользователя:**

```http
POST /users/setIsActive
Content-Type: application/json

{
  "user_id": "user-1",
  "is_active": false
}
```

**Получение PR’ов, назначенных пользователю:**

```http
GET /users/getReview?user_id=user-1
```

---

### Pull Requests

**Создание PR:**

```http
POST /pullRequest/create
Content-Type: application/json

{
  "pull_request_id": "pr-123",
  "pull_request_name": "Add new feature",
  "author_id": "user-1"
}
```

**Мерж PR:**

```http
POST /pullRequest/merge
Content-Type: application/json

{
  "pull_request_id": "pr-123"
}
```

**Reassign ревьюера:**

```http
POST /pullRequest/reassign
Content-Type: application/json

{
  "pull_request_id": "pr-123",
  "old_user_id": "user-2"
}
```

---

## Бизнес-логика

### 1. Назначение ревьюеров на новый PR

При создании PR сервис:

1. Находит команду, к которой относится автор PR.
2. Собирает список кандидатов-ревьюеров:

    * только активные пользователи (`is_active = true`);
    * исключается автор PR;
    * исключаются уже назначенные ревьюеры (при повторном вызове).
3. Перемешивает кандидатов с помощью **Fisher–Yates shuffle**.
4. Выбирает до двух ревьюеров (если людей меньше, назначает столько, сколько есть).
5. Сохраняет назначение в таблице `pr_reviewers`.

Это позволяет обеспечить справедливое и случайное распределение нагрузки.

### 2. Мерж PR

При запросе `/pullRequest/merge`:

* проверяется существование PR;
* проверяется, не замержен ли он уже:

    * если уже `MERGED` — операция идемпотентна, возвращаем текущий статус;
* обновляется статус PR, выставляется `merged_at`;
* возвращается актуальное состояние PR.

### 3. Reassign ревьюера

Запрос `/pullRequest/reassign`:

1. Проверяется, существует ли PR.
2. Проверяется, был ли на нём назначен `old_user_id`:

    * если нет — возвращается ошибка `NOT_ASSIGNED`.
3. Формируется список кандидатов:

    * только активные пользователи команды;
    * исключается автор PR;
    * исключается `old_user_id`;
    * исключаются уже назначенные ревьюеры на этот PR.
4. Если кандидатов нет — ошибка `NO_CANDIDATE`.
5. Иначе выбирается новый ревьюер (случайно), старый снимается, новый добавляется.

---

## Формат ошибок

Все ошибки возвращаются в JSON-формате:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "team not found"
  }
}
```

Сервис использует следующие коды ошибок:

* `TEAM_EXISTS` (400) — команда с таким именем уже существует;
* `USER_ALREADY_EXISTS` (409) — пользователь уже существует;
* `PR_EXISTS` (409) — PR уже существует;
* `NOT_FOUND` (404) — сущность не найдена (команда, пользователь, PR);
* `PR_MERGED` (409) — операция недопустима, PR уже замержен;
* `NOT_ASSIGNED` (409) — пользователь не был ревьюером данного PR;
* `NO_CANDIDATE` (409) — нет кандидатов для назначения ревьюера.

---

## Тестирование


### E2E-тесты

E2E-тесты прогоняют полный сценарий работы API на реальной PostgreSQL.

Примерный сценарий:

1. Подготовка окружения:

   ```bash
   make e2e-setup
   ```

2. Запуск сервиса для E2E:

   ```bash
   make e2e-run-api
   ```

3. Запуск E2E-тестов:

   ```bash
   make e2e-test
   ```

4. Остановка окружения:

   ```bash
   make e2e-teardown
   ```

Покрываются кейсы:

* создание PR и авто-назначение ревьюеров;
* merge PR (в т.ч. повторный вызов);
* reassign ревьюера;
* создание и получение команды;
* изменение активности пользователя;
* получение списка PR на ревью для пользователя.

---

## Makefile: полезные команды

```bash
# Миграции
make install-goose          # Установить goose
make migrate-up             # Применить миграции
make migrate-down           # Откатить миграции
make migrate-status         # Статус миграций
make migrate-create NAME=name  # Создать новый файл миграции

# Сборка и запуск
make build                  # Собрать бинарник
make run                    # Запустить приложение

# Тесты
make test                   # Unit-тесты
make test-coverage          # Тесты + отчёт покрытия
make e2e-setup              # Подготовка окружения E2E
make e2e-run-api            # Запуск API для E2E
make e2e-test               # E2E-тесты
make e2e-teardown           # Остановка E2E-окружения

# Качество кода
make lint                   # Линтер
make fmt                    # Форматирование кода

# Docker
make docker-up              # docker-compose up
make docker-down            # docker-compose down
make docker-logs            # Логи контейнеров

# Разное
make clean                  # Очистка артефактов
make help                   # Список команд
```

---

## Структура проекта

```text
.
├── cmd/
│   └── api/
│       └── main.go                 # Точка входа приложения
├── internal/
│   ├── app/
│   │   └── app.go                  # Инициализация и запуск приложения
│   ├── config/
│   │   └── config.go               # Загрузка конфигурации из env
│   ├── domain/
│   │   └── models/                 # Доменные сущности
│   │       ├── pr.go
│       ├── team.go
│       └── user.go
│   ├── dto/
│   │   ├── request/                # DTO запросов
│   │   │   ├── pr.go
│   │   │   ├── team.go
│   │   │   └── user.go
│   │   └── response/               # DTO ответов
│   │       ├── error.go
│   │       ├── pr.go
│   │       ├── team.go
│   │       └── user.go
│   ├── handler/                    # HTTP-обработчики
│   │   ├── health.go
│   │   ├── helpers.go
│   │   ├── pr.go
│   │   ├── team.go
│   │   └── user.go
│   ├── middleware/                 # HTTP-middleware
│   │   ├── logger.go
│   │   └── recovery.go
│   ├── repository/                 # Интерфейсы репозиториев и транзакций
│   │   ├── interfaces.go
│   │   └── postgres/
│   │       ├── helpers.go
│   │       ├── pr.go
│   │       ├── team.go
│   │       └── user.go
│   │   └── transaction.go
│   └── service/                    # Бизнес-логика
│       ├── interfaces.go
│       ├── pr.go
│       ├── team.go
│       └── user.go
├── pkg/
│   ├── database/
│   │   └── postgres.go             # Подключение к PostgreSQL
│   ├── errors/
│   │   └── errors.go               # Общий слой ошибок
│   └── logger/
│       └── logger.go               # Логирование
├── migrations/                     # SQL-миграции
│   ├── 00001_init_schema.sql
│   ├── 00002_create_teams.sql
│   ├── 00003_create_users.sql
│   ├── 00004_create_pull_requests.sql
│   └── 00005_create_pr_reviewers.sql
├── test/
│   └── e2e/                        # E2E-тесты
│       ├── main_test.go
│       ├── pr_test.go
│       ├── team_test.go
│       └── user_test.go
├── .gitignore
├── docker-compose.yml              # Production compose
├── docker-compose.e2e.yml          # Compose для E2E
├── Dockerfile                      # Multi-stage Dockerfile
├── entrypoint.sh                   # Entrypoint-скрипт
├── Makefile
├── README.md
└── TODO.md                         # План дальнейшей разработки
```

---

## Миграции:

1. `00001_init_schema.sql` — базовая схема;
2. `00002_create_teams.sql` — таблица `teams`;
3. `00003_create_users.sql` — таблица `users`;
4. `00004_create_pull_requests.sql` — таблица `pull_requests`;
5. `00005_create_pr_reviewers.sql` — таблица `pr_reviewers`.

---

## Пример сценария использования (через curl)

```bash
# 1. Создать команду
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend-team",
    "members": [
      {"user_id": "user-1", "username": "Alice", "is_active": true},
      {"user_id": "user-2", "username": "Bob", "is_active": true},
      {"user_id": "user-3", "username": "Charlie", "is_active": true}
    ]
  }'

# 2. Создать PR (автоматически назначатся до 2 ревьюеров)
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-123",
    "pull_request_name": "Add authentication",
    "author_id": "user-1"
  }'

# 3. Посмотреть PR’ы, которые должен ревьюить user-2
curl "http://localhost:8080/users/getReview?user_id=user-2"

# 4. Переназначить ревьюера
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-123",
    "old_user_id": "user-2"
  }'

# 5. Замержить PR
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-123"
  }'
```

---

## Заключение

Этот сервис реализует автоматическое назначение ревьюеров на Pull Request’ы и демонстрирует типичный подход к построению backend-сервиса с использованием Go, PostgreSQL, Docker и миграций базы данных.
