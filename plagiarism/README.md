# Plagiarism Service

Микросервис для асинхронной проверки сдач на плагиат. Принимает запрос на проверку, ставит задачу в очередь, воркер скачивает работы из filestorage и сравнивает байт‑к‑байту, складывая отчёты в файловую систему.

## Стек

- Go 1.25 (stdlib net/http, без внешних зависимостей)
- Внешний сервис filestorage (HTTP API)
- Локальное файловое хранилище отчётов (`plagiarism/reports`)
- Docker (многоэтапная сборка)

## Быстрый старт

```bash
cd plagiarism
# настроить адрес filestorage (должен быть доступен)
export FILESTORAGE_URL=http://localhost:8080
export PORT=8081

go run ./cmd/server
```

API слушает на `:$PORT` (по умолчанию 8081). Отчёты пишутся в `plagiarism/reports` относительно рабочего каталога.

## API

| Метод | Путь | Описание |
|-------|------|----------|
| `POST /checks` | JSON `{"submission_id": "...", "work_id": "..."}` | Ставит проверку в очередь, отвечает ACK `submission_id` + `status=pending`. |
| `GET /works/{work_id}/reports` | Возвращает последний известный отчёт по всем сдачам работы. |

Спека OpenAPI: `plagiarism/openapi.yaml`.

Пример (предполагается, что filestorage уже содержит сдачу с указанным `submission_id`):

```bash
curl -X POST http://localhost:8081/checks \
  -H "Content-Type: application/json" \
  -d '{"submission_id":"sub-1","work_id":"work-1"}'

curl "http://localhost:8081/works/work-1/reports" | jq .
```

В отчётах `matches` включают только совпадения выше порога `MATCH_THRESHOLD`.

## Переменные окружения

- `PORT` — порт HTTP сервера (по умолчанию `8081`).
- `FILESTORAGE_URL` — базовый URL filestorage (по умолчанию `http://localhost:8080`; важно указать реальный адрес, чтобы не ходить в себя).
- `MATCH_THRESHOLD` — порог совпадения, 0…1 (по умолчанию `0.8`).
- `WORKER_COUNT` — количество параллельных воркеров (по умолчанию `1`).

## Структура проекта

- `cmd/server/main.go` — точка входа, DI, конфигурация.
- `internal/api/http` — хендлеры и маршрутизация.
- `internal/application/usecase` — бизнес‑логика (старт проверки, получение отчётов).
- `internal/domain` — модели `CheckReport`, `MatchResult`.
- `internal/infrastructure` — адаптеры: конфиг, filestorage клиент, файловое хранилище отчётов, воркер.

## Docker

Сборка образа:

```bash
cd plagiarism
docker build -t plagiarism .
```

Запуск с внешним filestorage (пример адреса и проброса порта):

```bash
docker run --rm \
  -e PORT=8081 \
  -e FILESTORAGE_URL=http://filestorage:8080 \
  -e MATCH_THRESHOLD=0.8 \
  -e WORKER_COUNT=2 \
  -p 8081:8081 \
  -v "$(pwd)/reports:/app/plagiarism/reports" \
  plagiarism
```

Файлы отчётов сохраняются в `/app/plagiarism/reports` (смонтировано volume). Для полноценной работы нужен запущенный filestorage и наличие нужных `submission_id` в нём.
