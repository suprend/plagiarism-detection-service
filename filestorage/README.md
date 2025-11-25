# Filestorage Service

Микросервис отвечает за загрузку работ студентов, хранение метаданных в Postgres и файлов в S3‑совместимом хранилище (MinIO). API позволяет отправить решение, получить список сдач по заданию и скачать конкретный файл.

## Стек

- Go 1.25 (stdlib net/http, pgx, AWS SDK v2, sqlc)
- Postgres 16 (таблица `submissions`)
- MinIO (S3 совместимый storage)
- Docker + Docker Compose для локального окружения

## Быстрый старт

```bash
cd filestorage
docker compose up --build
```

Compose поднимет Postgres, MinIO, job инициализации бакета и само приложение (`app` на `localhost:8080`).

## API

| Метод | Путь | Описание |
|-------|------|----------|
| `POST /submit` | multipart form (`assignment_id`, `login`, `file`) | Создаёт submission и грузит файл в S3. Лимит размера — 1 МБ. |
| `GET /submissions?assignment_id=...` | Возвращает список сдач для задания. |
| `GET /submissions/download?submission_id=...` | Стримит файл по `submission_id`. Имя и тип в ответе — `submission_id` + `application/octet-stream`. |

Пример ручного теста (из корня репозитория с `tmp-files/sample1.txt`):

```bash
curl -X POST http://localhost:8080/submit \
  -F assignment_id=test-assignment \
  -F login=test-user \
  -F file=@tmp-files/sample1.txt

curl "http://localhost:8080/submissions?assignment_id=test-assignment" | jq .

curl -L "http://localhost:8080/submissions/download?submission_id=<ID>" \
  -o tmp-files/downloaded.bin
```

## Переменные окружения

В `docker-compose.yml` уже указаны дефолты:

- `DATABASE_URL` — строка подключения к Postgres (контейнер `postgres`).
- `S3_BUCKET`, `S3_ENDPOINT`, `AWS_*` — настройки MinIO.

При запуске вне Compose их нужно задать вручную. Ограничение размера файла настраивается в `internal/api/http/handler/submit.go` (константа `maxUploadSize`).

## Структура проекта

- `cmd/server/main.go` — точка входа, конфигурация, DI.
- `internal/api/http` — хендлеры, маршрутизация, API ошибки.
- `internal/application/usecase` — бизнес‑логика (submit / download / get submissions).
- `internal/domain` — сущности и интерфейсы репозиториев.
- `internal/infrastructure/repository/postgres` — sqlc‑генерированные запросы и адаптер.
- `internal/infrastructure/repository/s3` — работа с MinIO/S3.
- `migrations/` — SQL для таблицы `submissions`.

## Docker

- `Dockerfile` — двухэтапная сборка: кэш зависимостей, `GOARCH` через `TARGETARCH`, запуск от непривилегированного пользователя, OCI‑labels.
- `.dockerignore` — исключает git, IDE, временные каталоги (`tmp`, `tmp-files`, логи и т.д.).
- `docker-compose.yml` — запускает Postgres, MinIO, job и приложение с healthcheck’ами и `depends_on.condition`.

