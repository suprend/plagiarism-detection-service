## User API Gateway

Лёгкий шлюз, который принимает запросы клиентов, прокидывает их в `filestorage` и `plagiarism`, а затем возвращает агрегированный ответ.

### Быстрый старт

```bash
cd userapi
export FILESTORAGE_URL=http://localhost:8080    # адрес filestorage
export PLAGIARISM_URL=http://localhost:8081    # адрес plagiarism
export PORT=8082                               # userapi listen port

go run ./cmd/server
```

Ниже лежащие сервисы должны быть запущены и доступны по указанным адресам.

### Docker

```bash
cd userapi
docker build -t userapi .
docker run --rm -p 8082:8082 \
  -e FILESTORAGE_URL=http://filestorage:8080 \
  -e PLAGIARISM_URL=http://plagiarism:8081 \
  userapi
```

Swagger UI доступен на `http://localhost:8082/swagger`, сама спецификация — `http://localhost:8082/openapi.yaml`.

### API

- `POST /works/{work_id}/submit` — multipart с полями `login` (string) и `file` (<=1MB). Загружает решение в filestorage и сразу ставит задачу на проверку плагиата. Ответ: `{"submission_id":"...","check_status":"pending"}` с HTTP 202.
- `GET /works/{work_id}/reports` — проксирует последние отчёты по работе из сервиса plagiarism. Формат совпадает с его API (`{"work_id":"...","reports":[...]}`).
- `GET /wordcloud?submission_id=...` — строит облако слов для конкретной сдачи (png), использует сервис quickchart.io.

### Конфигурация

- `PORT` — порт HTTP (по умолчанию `8082`).
- `FILESTORAGE_URL` — базовый адрес filestorage (по умолчанию `http://localhost:8080`).
- `PLAGIARISM_URL` — базовый адрес plagiarism (по умолчанию `http://localhost:8081`).
- `MAX_UPLOAD_SIZE_BYTES` — лимит размера загружаемого файла (по умолчанию `1048576`, то есть 1MB).
- `WORDCLOUD_URL` — endpoint сервиса построения облака слов (по умолчанию `https://quickchart.io/wordcloud`).
- `WORDCLOUD_DIR` — путь для сохранения сгенерированных PNG облаков слов (по умолчанию `tmp-files/wordclouds` внутри контейнера/проекта).
