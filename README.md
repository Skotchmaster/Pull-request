# PR Reviewer Assignment Service

Сервис для автоматического назначения ревьюверов на Pull Request’ы внутри команд, а также управления командами и пользователями.

---

## Стек технологий

- Язык: Go  
- Web-фреймворк: `github.com/labstack/echo/v4`  
- ORM: `gorm.io/gorm` + `gorm.io/driver/postgres`  
- База данных: PostgreSQL  
- Контейнеризация: Docker + docker-compose  
- Логирование: простой обёрнутый логгер (`internal/logging`)

---

## Структура проекта

```text
cmd/server/main.go        – точка входа, конфиг, запуск сервера, graceful shutdown

internal/config           – загрузка конфигурации из env (DB_URL, PORT, JWT_SECRET)
internal/db               – инициализация подключения к БД (GORM)
internal/logging          – интерфейс и реализация логгера
internal/http/route.go    – регистрация HTTP-роутов в Echo
internal/models           – модели БД (Team, User, PullRequest, PRReviewer)
internal/repo             – слой доступа к БД (репозитории на GORM)
internal/service          – бизнес-логика (назначение ревьюверов, merge, переназначение)
internal/handlers         – HTTP-хэндлеры, маппинг доменных ошибок в HTTP-коды
internal/apierr           – доменные ошибки и формат ответа ErrorResponse

db/migrations             – SQL-миграции (инициализация схемы БД)
openapi.yml               – OpenAPI-спецификация сервиса
Dockerfile                – сборка контейнера приложения
docker-compose.yml        – запуск Postgres и сервиса
.env.example              – пример .env-файла
Makefile                  – вспомогательные команды (*nix)
```

---

## Запуск через Docker

Требуется установленный Docker и docker-compose.

1. Скопировать файл окружения:

```bash
cp .env.example .env
```

2. Заполнить значения в `.env`:

```env
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=pull_request

DB_URL=postgres://postgres:postgres@db:5432/pull_request?sslmode=disable

PORT=8080
```

3. Поднять сервис:

```bash
docker-compose up --build
```

или через Makefile:

```bash
make up
```

После запуска:

- HTTP API: `http://localhost:8080`
- Liveness-эндпоинт: `GET /health/live`

Миграции применяются автоматически при старте контейнера Postgres  
(через `./db/migrations` и `docker-entrypoint-initdb.d`).

---

## Makefile

В корне проекта есть `Makefile` с основными командами для разработки и запуска:

```makefile
BINARY_NAME = pr-service

.PHONY: build run test tidy up down logs

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

up:
	docker-compose up --build

down:
	docker-compose down -v

logs:
	docker-compose logs -f
```

Полезные команды:

- `make build` – собрать бинарник `bin/pr-service`.  
- `make run` – запустить сервис локально без Docker.  
- `make test` – прогнать все тесты `go test ./...`.  
- `make tidy` – привести зависимости в порядок (`go mod tidy`).  
- `make lint` – запустить линтер `golangci-lint run ./...`.  
- `make up` – поднять Postgres и сервис через `docker-compose up --build`.  
- `make down` – остановить и удалить контейнеры и volume’ы (`docker-compose down -v`).  
- `make logs` – посмотреть логи `docker-compose logs -f`.

---

## API и контракт

Полное описание контрактов – в `openapi.yml`.

Основные эндпоинты (см. `internal/http/route.go`):

- `GET /health/live` – liveness-проба.
- `POST /team/add` – создать команду с участниками (одновременно создаёт/обновляет пользователей).
- `GET /team/get` – получить информацию о команде и её участниках.
- `GET /users/getReview` – список PR’ов, где пользователь назначен ревьювером.
- `POST /users/setIsActive` – изменение флага активности пользователя.
- `POST /pullRequest/create` – создать PR и автоматически назначить до двух ревьюверов.
- `POST /pullRequest/merge` – пометить PR как `MERGED` (идемпотентно).
- `POST /pullRequest/reassign` – переназначить ревьювера на случайного активного участника его команды.
- `GET /stats/reviewers` – простая статистика: сколько раз каждый пользователь был назначен ревьювером (на основе таблицы `pr_reviewers`).

Формат ошибок:

```json
{
  "error": {
    "code": "CODE",
    "message": "message"
  }
}
```

Список возможных значений `error.code` описан в `components.schemas.ErrorResponse` в `openapi.yml`, например:  
`TEAM_EXISTS`, `PR_EXISTS`, `PR_MERGED`, `NOT_ASSIGNED`, `NO_CANDIDATE`, `NOT_FOUND`, `BAD_REQUEST`, `INTERNAL`.

---

## Линтер

В проекте используется [golangci-lint](https://golangci-lint.run/) с конфигурацией в файле `.golangci.yml` в корне репозитория.

Основные включённые линтеры:

- `govet` — встроенный анализатор Go, ищет подозрительные конструкции;  
- `staticcheck` — расширенный статический анализ кода;  
- `unused` — поиск неиспользуемого кода;  
- `errcheck` — проверка, что ошибки не игнорируются.

### Установка и запуск

Установка:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Запуск:

```bash
make lint
```

или:

```bash
golangci-lint run ./...
```

---

## Обработка ошибок и принятые решения

### 1. 400 BAD_REQUEST и 500 INTERNAL в OpenAPI

В исходной постановке не были явно описаны:

- ответы `400` на некорректные запросы (битый JSON, отсутствие обязательных параметров и т.п.);
- ответы `500` при внутренних ошибках сервера.

На практике:

- запросы могут приходить с неполными или некорректными данными;
- возможны неожиданные ошибки БД или инфраструктуры.

Поэтому OpenAPI-спецификация была расширена:

- добавлен общий тип `ErrorResponse` с полем `code`, в котором есть в т.ч.:
  - `BAD_REQUEST` – ошибки валидации входных данных;
  - `INTERNAL` – необработанные внутренние ошибки;
- для большинства эндпоинтов добавлены ответы:
  - `400` – «Неверный запрос» / «Некорректные входные данные»;
  - `500` – «Внутренняя ошибка сервера».

Поведение в коде:

- `400 BAD_REQUEST` возвращается, если:
  - не хватает обязательного параметра (например, `user_id`, `team_name`);
  - тело запроса невалидно (ошибка при разборе JSON);
- `500 INTERNAL` возвращается, если:
  - происходит любая неожиданная ошибка (ошибка БД, паника и т.п.);
  - ошибка логируется на сервере, а клиент получает безопасное сообщение.

Допущение:  
Можно было бы считать, что все запросы заведомо валидны, но для реалистичного сервиса поведение при нештатных ситуациях должно быть явно определено и задокументировано.

---

### 2. Поведение `/users/getReview`, если `user_id` не существует в БД

В исходном задании и OpenAPI не было описано, что делать, если:

- клиент вызывает `GET /users/getReview?user_id=...`,
- а такого пользователя ещё нет в базе.

Возможные варианты:

1. Всегда возвращать `200` с пустым списком, независимо от существования пользователя.  
2. Возвращать `404`, если пользователь не найден.  
3. Возвращать `400` и считать это ошибкой клиента.

Выбран вариант №2 — возвращать `404 NOT_FOUND`, потому что:

- запрос звучит как «какие PR’ы у данного пользователя?»;
- если пользователя нет, это ближе к «ресурс не найден», чем к «просто нет PR’ов»;
- такое поведение отражено в `openapi.yml` для `/users/getReview`:
  - `404` с `error.code: NOT_FOUND` и сообщением `user not found`.

Фактическое поведение эндпоинта `/users/getReview`:

- если `user_id` не передан или пустой:
  - возвращается `400 BAD_REQUEST` с `code: BAD_REQUEST`;
- если пользователя с таким `user_id` нет в БД:
  - возвращается `404 NOT_FOUND` с `code: NOT_FOUND`;
- если пользователь существует, но не назначен ни на один PR:
  - возвращается `200 OK` с пустым массивом `pull_requests`.

---

### 3. Дополнительный эндпоинт статистики

В рамках дополнительных заданий был реализован простой статистический эндпоинт:

- `GET /stats/reviewers` — возвращает количество назначений на ревью по каждому пользователю.

Решение:

- статистика считается на лету на основе таблицы `pr_reviewers`;
- используется агрегирующий запрос `GROUP BY user_id`, без отдельной таблицы для статистики;
