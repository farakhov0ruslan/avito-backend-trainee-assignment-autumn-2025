# TODO: PR Reviewer Assignment Service

Проект разработан по Clean Architecture с использованием стандартной структуры Go проектов.

## Статус обозначений
- [ ] Не начато
- [x] Завершено
- [>] В процессе

---

## Фаза 0: Подготовка инфраструктуры

### 0.1. Инициализация зависимостей
- [x] Обновить go.mod с основными зависимостями (gorilla/mux, pgx/v5, godotenv)
- [x] Добавить зависимость для работы с PostgreSQL (pgx/v5, pgxpool)
- [x] Добавить зависимость для конфигурации (godotenv)
- [x] Установить goose для миграций (`go install github.com/pressly/goose/v3/cmd/goose@latest`)

### 0.2. Конфигурация окружения
- [x] Создать .env.example с примерами переменных (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, SERVER_PORT)
- [x] Создать .env для локальной разработки (не коммитить)

### 0.3. Makefile
- [x] Добавить команду `make install-goose` для установки goose
- [x] Добавить команду `make migrate-up` для применения миграций
- [x] Добавить команду `make migrate-down` для отката миграций
- [x] Добавить команду `make migrate-create NAME=<name>` для создания новой миграции
- [x] Добавить команду `make build` для сборки проекта
- [x] Добавить команду `make run` для запуска приложения
- [x] Добавить команду `make test` для запуска тестов
- [x] Добавить команду `make lint` для проверки кода линтером
- [x] Добавить команду `make docker-up` для запуска docker-compose
- [x] Добавить команду `make docker-down` для остановки docker-compose

---

## Фаза 1: Domain Models

**Важно:** Domain Models - это чистые структуры данных БЕЗ валидации!
- **Handler**: базовая валидация (парсинг JSON, типы данных)
- **Service**: бизнес-валидация (бизнес-правила)
- **Models**: только структуры данных + вспомогательные методы

### 1.1. User Domain Model (internal/domain/models/user.go)
- [x] Создать структуру User с полями: ID, Username, TeamName, IsActive
- [x] Добавить метод String() для удобного логирования
- [x] Добавить JSON и DB теги

### 1.2. Team Domain Model (internal/domain/models/team.go)
- [x] Создать структуру Team с полями: Name, Members ([]User)
- [x] Добавить метод GetActiveMembers() []User
- [x] Добавить метод GetActiveMembersExcept(userID string) []User
- [x] Добавить метод String() для логирования

### 1.3. PullRequest Domain Model (internal/domain/models/pr.go)
- [x] Создать enum для статуса PR (OPEN, MERGED)
- [x] Создать структуру PullRequest с полями: ID, Name, AuthorID, Status, AssignedReviewers, CreatedAt, MergedAt
- [x] Добавить метод IsMerged() bool
- [x] Добавить метод IsReviewerAssigned(userID string) bool
- [x] Добавить метод String() для логирования

---

## Фаза 2: DTO (Data Transfer Objects)

### 2.1. Request DTOs
- [x] internal/dto/request/team.go - CreateTeamRequest, GetTeamRequest
- [x] internal/dto/request/user.go - SetUserActiveRequest, GetUserReviewsRequest
- [x] internal/dto/request/pr.go - CreatePRRequest, MergePRRequest, ReassignReviewerRequest

### 2.2. Response DTOs
- [x] internal/dto/response/team.go - TeamResponse, CreateTeamResponse
- [x] internal/dto/response/user.go - UserResponse, SetUserActiveResponse, GetUserReviewsResponse
- [x] internal/dto/response/pr.go - PullRequestResponse, PullRequestShortResponse
- [x] internal/dto/response/pr.go - CreatePRResponse, MergePRResponse, ReassignReviewerResponse
- [x] internal/dto/response/error.go - ErrorResponse с code и message
- [x] internal/dto/response/error.go - Константы для error codes (TEAM_EXISTS, PR_EXISTS, PR_MERGED, NOT_ASSIGNED, NO_CANDIDATE, NOT_FOUND)

---

## Фаза 3: Утилиты и конфигурация

### 3.1. Custom Errors (pkg/errors/errors.go)
- [x] Создать кастомные ошибки: ErrTeamExists, ErrTeamNotFound
- [x] Создать кастомные ошибки: ErrUserNotFound, ErrUserAlreadyExists
- [x] Создать кастомные ошибки: ErrPRExists, ErrPRNotFound, ErrPRMerged
- [x] Создать кастомные ошибки: ErrReviewerNotAssigned, ErrNoCandidates
- [x] Добавить функции MapErrorToHTTPStatus и MapErrorToErrorCode

### 3.2. Logger (pkg/logger/logger.go)
- [x] Настроить структурированный логгер (стандартный log)
- [x] Добавить уровни логирования (DEBUG, INFO, WARN, ERROR)
- [x] Добавить методы Info, Error, Debug, Warn, Fatal

### 3.3. Database Connection (pkg/database/postgres.go)
- [x] Создать функцию NewPostgresDB для подключения через pgxpool
- [x] Добавить пинг базы данных для проверки подключения
- [x] Добавить настройку connection pool (max/min conns, lifetimes)
- [x] Добавить функцию Close для graceful shutdown

### 3.4. Config (internal/config/config.go)
- [x] Создать структуру Config с полями для БД, сервера, приложения
- [x] Создать функцию Load для загрузки из переменных окружения
- [x] Добавить валидацию обязательных параметров (Validate)
- [x] Добавить поддержку .env файла через godotenv

---

## Фаза 4: Миграции базы данных

### 4.1. Миграция: Init Schema (migrations/00001_init_schema.sql)
- [x] Создать up-миграцию (базовая, расширения не нужны)
- [x] Создать down-миграцию для отката

### 4.2. Миграция: Teams Table (migrations/00002_create_teams.sql)
- [x] Создать up-миграцию с таблицей teams (name PRIMARY KEY)
- [x] Добавить индексы (created_at)
- [x] Создать down-миграцию (DROP TABLE teams CASCADE)

### 4.3. Миграция: Users Table (migrations/00003_create_users.sql)
- [x] Создать up-миграцию с таблицей users (id PRIMARY KEY, username, team_name, is_active)
- [x] Добавить FOREIGN KEY на teams.name (ON DELETE CASCADE)
- [x] Добавить индексы (team_name, is_active, composite)
- [x] Создать down-миграцию (DROP TABLE users CASCADE)

### 4.4. Миграция: Pull Requests Table (migrations/00004_create_pull_requests.sql)
- [x] Создать up-миграцию с таблицей pull_requests (id PRIMARY KEY, name, author_id, status, created_at, merged_at)
- [x] Добавить FOREIGN KEY на users.id для author_id (ON DELETE CASCADE)
- [x] Добавить CHECK constraint для status (OPEN или MERGED)
- [x] Создать down-миграцию (DROP TABLE pull_requests CASCADE)

### 4.5. Миграция: PR Reviewers Junction Table (migrations/00005_create_pr_reviewers.sql)
- [x] Создать up-миграцию с таблицей pr_reviewers (pr_id, reviewer_id, assigned_at)
- [x] Добавить FOREIGN KEY на pull_requests.id и users.id (ON DELETE CASCADE)
- [x] Добавить PRIMARY KEY (pr_id, reviewer_id) - обеспечивает уникальность
- [x] Создать down-миграцию (DROP TABLE pr_reviewers CASCADE)

---

## Фаза 5: Repository Layer

### 5.1. Repository Interfaces (internal/repository/interfaces.go)
- [x] Создать интерфейс TeamRepository с методами Create, GetByName, Exists
- [x] Создать интерфейс UserRepository с методами Create, Update, GetByID, GetByTeamName, SetActive
- [x] Создать интерфейс PRRepository с методами Create, GetByID, Update, Merge, GetReviewersByPRID, AddReviewer, RemoveReviewer, GetPRsByReviewerID

### 5.2. Transaction Manager (internal/repository/transaction.go)
- [x] Создать интерфейс TransactionManager с методом WithTransaction
- [x] Реализовать PgxTransactionManager с поддержкой транзакций
- [x] Добавить Executor интерфейс и функцию GetTx для работы с транзакциями

### 5.3. Team Repository (internal/repository/postgres/team.go)
- [x] Реализовать метод Create для создания команды
- [x] Реализовать метод GetByName для получения команды с участниками
- [x] Реализовать метод Exists для проверки существования команды
- [x] Добавить обработку ошибок (unique violation) и логирование

### 5.4. User Repository (internal/repository/postgres/user.go)
- [ ] Реализовать метод Create для создания пользователя
- [ ] Реализовать метод Update для обновления пользователя (upsert логика)
- [ ] Реализовать метод GetByID для получения пользователя по ID
- [ ] Реализовать метод GetByTeamName для получения всех пользователей команды
- [ ] Реализовать метод SetActive для изменения флага is_active
- [ ] Добавить обработку ошибок и логирование

### 5.5. PR Repository (internal/repository/postgres/pr.go)
- [ ] Реализовать метод Create для создания PR
- [ ] Реализовать метод GetByID для получения PR по ID (с ревьюверами)
- [ ] Реализовать метод Update для обновления PR
- [ ] Реализовать метод Merge для изменения статуса на MERGED (идемпотентно)
- [ ] Реализовать метод GetReviewersByPRID для получения списка ревьюверов
- [ ] Реализовать метод AddReviewer для добавления ревьювера
- [ ] Реализовать метод RemoveReviewer для удаления ревьювера
- [ ] Реализовать метод GetPRsByReviewerID для получения PR'ов где пользователь ревьювер
- [ ] Добавить обработку ошибок и логирование

---

## Фаза 6: Service Layer (Бизнес-логика)

### 6.1. Service Interfaces (internal/service/interfaces.go)
- [ ] Создать интерфейс TeamService с методами CreateTeam, GetTeam
- [ ] Создать интерфейс UserService с методами SetUserActive, GetUserReviews
- [ ] Создать интерфейс PRService с методами CreatePR, MergePR, ReassignReviewer

### 6.2. Team Service (internal/service/team.go)
- [ ] Реализовать метод CreateTeam (проверка на существование + создание команды и пользователей)
- [ ] Добавить транзакцию для атомарного создания команды и пользователей
- [ ] Реализовать метод GetTeam (получение команды с участниками)
- [ ] Добавить обработку ошибок и валидацию

### 6.3. User Service (internal/service/user.go)
- [ ] Реализовать метод SetUserActive (изменение флага активности)
- [ ] Реализовать метод GetUserReviews (получение PR'ов где пользователь ревьювер)
- [ ] Добавить обработку ошибок и валидацию

### 6.4. PR Service - Часть 1: Create (internal/service/pr.go)
- [ ] Реализовать метод CreatePR
- [ ] Добавить проверку существования автора
- [ ] Добавить проверку существования PR (должен быть уникальным)
- [ ] Получить команду автора
- [ ] Реализовать логику выбора до 2 активных ревьюверов (исключая автора)
- [ ] Использовать случайный выбор ревьюверов
- [ ] Добавить транзакцию для создания PR и назначения ревьюверов
- [ ] Обработать случай когда доступных ревьюверов меньше 2

### 6.5. PR Service - Часть 2: Merge (internal/service/pr.go)
- [ ] Реализовать метод MergePR
- [ ] Проверить существование PR
  - [ ] Проверить что PR еще не MERGED (если уже MERGED - вернуть текущее состояние, идемпотентность)
- [ ] Изменить статус на MERGED и установить merged_at
- [ ] Добавить обработку ошибок

### 6.6. PR Service - Часть 3: Reassign (internal/service/pr.go)
- [ ] Реализовать метод ReassignReviewer
- [ ] Проверить существование PR
- [ ] Проверить что PR не MERGED (если MERGED - вернуть ошибку PR_MERGED)
- [ ] Проверить что old_user_id назначен ревьювером (если нет - ошибка NOT_ASSIGNED)
- [ ] Получить команду заменяемого пользователя
- [ ] Получить список активных кандидатов из команды (исключая old_user_id и автора PR)
- [ ] Выбрать случайного кандидата
- [ ] Если кандидатов нет - вернуть ошибку NO_CANDIDATE
- [ ] Использовать транзакцию для удаления старого и добавления нового ревьювера
- [ ] Добавить обработку ошибок

---

## Фаза 7: Handler Layer (HTTP API)

### 7.1. Health Handler (internal/handler/health.go)
- [ ] Создать структуру HealthHandler
- [ ] Реализовать метод Check для GET /health (вернуть 200 OK)
- [ ] Добавить проверку подключения к БД (опционально)

### 7.2. Team Handlers (internal/handler/team.go)
- [ ] Создать структуру TeamHandler с зависимостью от TeamService
- [ ] Реализовать метод CreateTeam для POST /team/add
  - [ ] Парсинг JSON из request body
  - [ ] Валидация входных данных
  - [ ] Вызов service.CreateTeam
  - [ ] Обработка ошибок (400 если TEAM_EXISTS, 500 для остальных)
  - [ ] Формирование response (201 Created)
- [ ] Реализовать метод GetTeam для GET /team/get
  - [ ] Получение team_name из query параметра
  - [ ] Вызов service.GetTeam
  - [ ] Обработка ошибок (404 если NOT_FOUND)
  - [ ] Формирование response (200 OK)

### 7.3. User Handlers (internal/handler/user.go)
- [ ] Создать структуру UserHandler с зависимостью от UserService
- [ ] Реализовать метод SetIsActive для POST /users/setIsActive
  - [ ] Парсинг JSON из request body
  - [ ] Валидация входных данных
  - [ ] Вызов service.SetUserActive
  - [ ] Обработка ошибок (404 если NOT_FOUND)
  - [ ] Формирование response (200 OK)
- [ ] Реализовать метод GetUserReviews для GET /users/getReview
  - [ ] Получение user_id из query параметра
  - [ ] Вызов service.GetUserReviews
  - [ ] Обработка ошибок (404 если NOT_FOUND)
  - [ ] Формирование response (200 OK)

### 7.4. PR Handlers (internal/handler/pr.go)
- [ ] Создать структуру PRHandler с зависимостью от PRService
- [ ] Реализовать метод CreatePR для POST /pullRequest/create
  - [ ] Парсинг JSON из request body
  - [ ] Валидация входных данных
  - [ ] Вызов service.CreatePR
  - [ ] Обработка ошибок (404 если NOT_FOUND, 409 если PR_EXISTS)
  - [ ] Формирование response (201 Created)
- [ ] Реализовать метод MergePR для POST /pullRequest/merge
  - [ ] Парсинг JSON из request body
  - [ ] Валидация входных данных
  - [ ] Вызов service.MergePR
  - [ ] Обработка ошибок (404 если NOT_FOUND)
  - [ ] Формирование response (200 OK, идемпотентно)
- [ ] Реализовать метод ReassignReviewer для POST /pullRequest/reassign
  - [ ] Парсинг JSON из request body
  - [ ] Валидация входных данных
  - [ ] Вызов service.ReassignReviewer
  - [ ] Обработка ошибок (404 если NOT_FOUND, 409 для PR_MERGED/NOT_ASSIGNED/NO_CANDIDATE)
  - [ ] Формирование response (200 OK)

---

## Фаза 8: Middleware

### 8.1. Logger Middleware (internal/middleware/logger.go)
- [ ] Создать middleware для логирования HTTP запросов
- [ ] Логировать: метод, путь, status code, время выполнения
- [ ] Использовать настроенный logger из pkg/logger

### 8.2. Recovery Middleware (internal/middleware/recovery.go)
- [ ] Создать middleware для восстановления от panic
- [ ] Логировать stack trace при panic
- [ ] Возвращать 500 Internal Server Error при panic

---

## Фаза 9: Application Initialization

### 9.1. Router Setup (internal/app/app.go)
- [ ] Создать функцию NewRouter с настройкой gorilla/mux
- [ ] Зарегистрировать health endpoint (GET /health)
- [ ] Зарегистрировать team endpoints (POST /team/add, GET /team/get)
- [ ] Зарегистрировать user endpoints (POST /users/setIsActive, GET /users/getReview)
- [ ] Зарегистрировать PR endpoints (POST /pullRequest/create, POST /pullRequest/merge, POST /pullRequest/reassign)
- [ ] Подключить middleware (logger, recovery)

### 9.2. Dependency Injection (internal/app/app.go)
- [ ] Создать функцию NewApp для инициализации всех зависимостей
- [ ] Инициализировать logger
- [ ] Инициализировать database connection
- [ ] Инициализировать repositories
- [ ] Инициализировать services
- [ ] Инициализировать handlers
- [ ] Инициализировать router

### 9.3. App Run (internal/app/app.go)
- [ ] Создать метод Run для запуска HTTP сервера
- [ ] Добавить graceful shutdown (обработка SIGINT, SIGTERM)
- [ ] Добавить таймауты для сервера (read, write, idle timeout)

---

## Фаза 10: Main Entry Point

### 10.1. Main (cmd/api/main.go)
- [ ] Загрузить конфигурацию
- [ ] Инициализировать приложение через app.NewApp
- [ ] Запустить сервер через app.Run
- [ ] Добавить обработку ошибок и логирование

---

## Фаза 11: Docker и Deployment

### 11.1. Dockerfile
- [ ] Использовать multi-stage build (golang:1.24 для сборки, alpine для runtime)
- [ ] Скопировать исходный код
- [ ] Собрать бинарник (go build -o /app cmd/api/main.go)
- [ ] В runtime образе: скопировать бинарник и миграции
- [ ] Установить goose в runtime образе для применения миграций
- [ ] EXPOSE 8080
- [ ] Создать entrypoint скрипт для применения миграций перед запуском

### 11.2. Entrypoint Script
- [ ] Создать entrypoint.sh скрипт
- [ ] Добавить команду применения миграций (goose -dir /migrations postgres "connection_string" up)
- [ ] Добавить запуск приложения
- [ ] Сделать скрипт исполняемым

### 11.3. docker-compose.yml
- [ ] Настроить сервис postgres (образ postgres:15, environment variables, volume для данных)
- [ ] Добавить healthcheck для postgres
- [ ] Настроить сервис api (build из Dockerfile, depends_on postgres с condition: service_healthy)
- [ ] Прокинуть порт 8080
- [ ] Настроить переменные окружения через env_file или environment
- [ ] Добавить networks для изоляции

### 11.4. Тестирование запуска
- [ ] Запустить docker-compose up
- [ ] Проверить что БД поднялась
- [ ] Проверить что миграции применились
- [ ] Проверить что сервис доступен на localhost:8080
- [ ] Проверить endpoint /health
- [ ] Исправить возможные ошибки

---

## Фаза 12: E2E Тестирование

### 12.1. docker-compose.e2e.yml
- [ ] Создать отдельный docker-compose для e2e тестов
- [ ] Настроить отдельную БД для тестов (db_e2e)
- [ ] Настроить сервис api_e2e с env файлом .env.e2e
- [ ] Создать сервис для запуска тестов (tests)
- [ ] Добавить depends_on с правильным порядком запуска

### 12.2. E2E тесты для Team (tests/e2e/team_test.go)
- [ ] Написать тест для создания команды (POST /team/add)
- [ ] Написать тест для получения команды (GET /team/get)
- [ ] Написать тест для ошибки TEAM_EXISTS
- [ ] Написать тест для ошибки NOT_FOUND

### 12.3. E2E тесты для User (tests/e2e/user_test.go)
- [ ] Написать тест для установки is_active (POST /users/setIsActive)
- [ ] Написать тест для получения PR'ов пользователя (GET /users/getReview)
- [ ] Написать тест для ошибки NOT_FOUND

### 12.4. E2E тесты для PR (tests/e2e/pr_test.go)
- [ ] Написать тест для создания PR с автоназначением ревьюверов (POST /pullRequest/create)
- [ ] Написать тест для создания PR когда доступен только 1 ревьювер
- [ ] Написать тест для создания PR когда нет доступных ревьюверов
- [ ] Написать тест для merge PR (POST /pullRequest/merge)
- [ ] Написать тест для идемпотентности merge (повторный вызов)
- [ ] Написать тест для reassign ревьювера (POST /pullRequest/reassign)
- [ ] Написать тест для ошибки PR_MERGED при reassign
- [ ] Написать тест для ошибки NOT_ASSIGNED при reassign
- [ ] Написать тест для ошибки NO_CANDIDATE при reassign
- [ ] Написать тест для ошибки PR_EXISTS

### 12.5. Запуск E2E тестов
- [ ] Создать команду в Makefile для запуска e2e тестов
- [ ] Запустить тесты и проверить что все проходят
- [ ] Исправить найденные баги

---

## Фаза 13: Нагрузочное тестирование

### 13.1. K6 скрипт (tests/load/k6_script.js)
- [ ] Установить k6 (https://k6.io/docs/getting-started/installation/)
- [ ] Написать сценарий для создания команд
- [ ] Написать сценарий для создания пользователей
- [ ] Написать сценарий для создания PR'ов
- [ ] Написать сценарий для reassign ревьюверов
- [ ] Настроить параметры нагрузки (VUs, duration)
- [ ] Добавить checks для проверки корректности ответов
- [ ] Добавить thresholds для метрик (http_req_duration < 300ms, http_req_failed < 0.1%)

### 13.2. Запуск нагрузочных тестов
- [ ] Запустить k6 run tests/load/k6_script.js
- [ ] Собрать метрики (response time, throughput, error rate)
- [ ] Задокументировать результаты в README.md
- [ ] При необходимости оптимизировать узкие места

---

## Фаза 14: Документация и финализация

### 14.1. README.md
- [ ] Описать проект и его назначение
- [ ] Добавить инструкции по запуску (docker-compose up)
- [ ] Описать API endpoints (можно сослаться на openapi.yml)
- [ ] Описать структуру проекта
- [ ] Добавить раздел с допущениями и решениями
  - [ ] Как реализован случайный выбор ревьюверов
  - [ ] Как обеспечивается идемпотентность merge
  - [ ] Как работает переназначение (из команды заменяемого)
- [ ] Добавить раздел с результатами нагрузочного тестирования
- [ ] Добавить примеры запросов (curl или httpie)

### 14.2. Линтер конфигурация
- [ ] Создать .golangci.yml с настройками линтера
- [ ] Включить основные линтеры (govet, errcheck, staticcheck, gosimple, ineffassign, unused)
- [ ] Настроить правила для импортов (gci)
- [ ] Добавить команду make lint в Makefile

### 14.3. Финальная проверка
- [ ] Убедиться что docker-compose up работает из коробки
- [ ] Убедиться что сервис доступен на порту 8080
- [ ] Проверить все endpoints через curl/Postman
- [ ] Запустить все тесты
- [ ] Запустить линтер
- [ ] Проверить что .env не закоммичен
- [ ] Убедиться что README.md полный и понятный

---

## Дополнительные задачи (бонус)

### Бонус 1: Эндпоинт статистики
- [ ] Добавить GET /stats endpoint
- [ ] Показать количество PR'ов по статусам
- [ ] Показать количество назначений по пользователям
- [ ] Добавить repository метод для получения статистики
- [ ] Добавить service метод
- [ ] Добавить handler

### Бонус 2: Массовая деактивация команды
- [ ] Добавить POST /team/deactivate endpoint
- [ ] Принимает team_name в body
- [ ] Деактивирует всех пользователей команды
- [ ] Переназначает открытые PR'ы где они были ревьюверами
- [ ] Оптимизировать для выполнения за ~100ms
- [ ] Использовать bulk операции в БД
- [ ] Добавить тесты

---

## Итоговая структура проекта

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   └── models/
│   │       ├── pr.go
│   │       ├── team.go
│   │       └── user.go
│   ├── dto/
│   │   ├── request/
│   │   │   ├── pr.go
│   │   │   ├── team.go
│   │   │   └── user.go
│   │   └── response/
│   │       ├── error.go
│   │       ├── pr.go
│   │       ├── team.go
│   │       └── user.go
│   ├── handler/
│   │   ├── health.go
│   │   ├── pr.go
│   │   ├── team.go
│   │   └── user.go
│   ├── middleware/
│   │   ├── logger.go
│   │   └── recovery.go
│   ├── repository/
│   │   ├── interfaces.go
│   │   ├── postgres/
│   │   │   ├── pr.go
│   │   │   ├── team.go
│   │   │   └── user.go
│   │   └── transaction.go
│   └── service/
│       ├── interfaces.go
│       ├── pr.go
│       ├── team.go
│       └── user.go
├── pkg/
│   ├── database/
│   │   └── postgres.go
│   ├── errors/
│   │   └── errors.go
│   └── logger/
│       └── logger.go
├── migrations/
│   ├── 00001_init_schema.sql
│   ├── 00002_create_teams.sql
│   ├── 00003_create_users.sql
│   ├── 00004_create_pull_requests.sql
│   └── 00005_create_pr_reviewers.sql
├── tests/
│   ├── e2e/
│   │   ├── pr_test.go
│   │   ├── team_test.go
│   │   └── user_test.go
│   └── load/
│       └── k6_script.js
├── .env.example
├── .gitignore
├── .golangci.yml
├── docker-compose.e2e.yml
├── docker-compose.yml
├── Dockerfile
├── entrypoint.sh
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── TODO.md
```

---

## Примечания

- Следовать принципам Clean Architecture
- Слабая связанность: использовать интерфейсы
- Handlers только парсят/валидируют, вся логика в Service
- Repository работает только с БД
- Использовать транзакции для атомарных операций
- Обрабатывать все ошибки
- Логировать важные события
- Тестировать критичную бизнес-логику
