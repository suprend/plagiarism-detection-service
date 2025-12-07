## Описание

Микросервисная система «антиплагиат»: студенты загружают работы, сервис их сохраняет, сравнивает между собой и выдаёт отчёты с совпадениями. Внешний API — только через `userapi` (gateway); хранение файлов и метаданных — в `filestorage`; анализ и отчёты — в `plagiarism`; построение облака слов — в отдельном `wordcloud`-сервисе.

## Архитектура

```mermaid
flowchart LR
    Client[[Клиент]]
    UserAPI[(userapi<br/>API Gateway)]
    FS[(filestorage)]
    PL[(plagiarism)]
    WC[(wordcloud)]
    QQ[(QuickChart API)]

    Client -- POST /works/{id}/submit --> UserAPI
    UserAPI -- POST /submit --> FS
    UserAPI -- POST /checks --> PL

    Client -- GET /works/{id}/reports --> UserAPI
    UserAPI -- GET /works/{id}/reports --> PL
    PL -- list/download submissions --> FS
    UserAPI -- JSON отчётов (с author_id) --> Client

    Client -- GET /wordcloud?submission_id=... --> UserAPI
    UserAPI -- GET /wordcloud --> WC
    WC -- download submission --> FS
    WC -- render cloud --> QQ
    UserAPI -- PNG --> Client
```

Микросервисы:
- `filestorage` — upload/list/download, метаданные в Postgres, файлы в MinIO.
- `plagiarism` — очередь проверок, воркер сравнивает сдачи, сохраняет отчёты (с author_id и other_author_id).
- `wordcloud` — строит облака слов на базе QuickChart, скачивая текст из filestorage.
- `userapi` — REST-шлюз: submit, reports, wordcloud. Swagger UI на `/swagger`.

## Алгоритм проверки плагиата

1. `userapi` после загрузки ставит задачу в `plagiarism` (`/checks`), статус сразу `pending`.
2. Воркер `plagiarism` получает все сдачи нужной работы из `filestorage` (`/submissions?assignment_id=...`), скачивает текущую и каждую чужую.
3. Сравнение — побайтово: считаем долю совпавших байт относительно большей длины двух файлов `similarity = matchedBytes / max(len(A), len(B))`.
4. Если `similarity >= MATCH_THRESHOLD` (по умолчанию 0.8), фиксируем совпадение с указанием `other_submission_id` и `other_author_id`.
5. По итогам пишется отчёт: `status=done` с найденными совпадениями или `failed` при ошибке скачивания/очереди; отчёты лежат в `plagiarism/reports/{work_id}/{submission_id}.json`, агрегат `overall.json`.

## Структура репозитория

- `filestorage/` — сервис хранения (cmd, api/http, usecase, infra, migrations).
- `plagiarism/` — сервис анализа (cmd, api/http, usecase, infra, воркер).
- `wordcloud/` — сервис построения облаков слов.
- `userapi/` — gateway (cmd, api/http, usecase, infra).
- `userapi/openapi.yaml` — OpenAPI спека публичного API.
- `filestorage/openapi.yaml` — OpenAPI спека сервиса хранения.
- `plagiarism/openapi.yaml` — OpenAPI спека сервиса проверки.
- `wordcloud/openapi.yaml` — OpenAPI спека сервиса облака слов.
- `docker-compose.yml` — общий запуск всех сервисов + Postgres + MinIO.
- `docs/` — вспомогательные материалы.

## Запуск

```bash
docker compose up --build
```

Поднимутся:
- filestorage: `localhost:8080`
- plagiarism: `localhost:8081`
- userapi: `localhost:8082`
- wordcloud: `localhost:8083`

Volumes: `postgres_data`, `minio_data`, `plagiarism_reports`. PNG wordcloud — в `tmp-files/wordclouds` внутри wordcloud-сервиса (можно примонтировать).

## Быстрый тест (curl)

```bash
# отправить работу
curl -X POST http://localhost:8082/works/test-work/submit \
  -F login=test-user \
  -F file=@tmp-files/icecream.txt

# получить отчёты
curl http://localhost:8082/works/test-work/reports | jq .

# облако слов (подставь submission_id из submit)
SID=<submission_id>
curl -o wordcloud.png "http://localhost:8082/wordcloud?submission_id=$SID"
```

Swagger UI: `http://localhost:8082/swagger`, спека: `/openapi.yaml`.

## Конфигурация (основные env)

- `MAX_UPLOAD_SIZE_BYTES` — лимит загрузки (filestorage/userapi).
- `MATCH_THRESHOLD`, `WORKER_COUNT` — plagiarism.
- `PORT`, `FILESTORAGE_URL`, `PLAGIARISM_URL`, `WORDCLOUD_SERVICE_URL` — адреса и порты сервисов.
- `WORDCLOUD_GENERATOR_URL`, `WORDCLOUD_DIR` — настройки сервиса wordcloud (по умолчанию QuickChart + `tmp-files/wordclouds`).
