## Wordcloud Service

Сервис строит облака слов по содержимому сдач, загруженных в File Storage.

- `GET /wordcloud?submission_id=...` — возвращает PNG. Если текст пустой или QuickChart недоступен, ответ будет с ошибкой 502.
- ENV:
  - `PORT` — порт HTTP (по умолчанию 8083).
  - `FILESTORAGE_URL` — базовый URL File Storage (по умолчанию `http://localhost:8080`).
  - `WORDCLOUD_GENERATOR_URL` — endpoint QuickChart (по умолчанию `https://quickchart.io/wordcloud`).
  - `WORDCLOUD_DIR` — путь для сохранения PNG (по умолчанию `tmp-files/wordclouds` внутри контейнера/проекта).

### Быстрый запуск локально

```bash
go run ./cmd/server
```

Далее запрос:

```bash
curl -o wordcloud.png "http://localhost:8083/wordcloud?submission_id=XXX"
```
