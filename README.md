# pz3-http — минимальный CRUD на Go (stdlib net/http)

Учебный проект: реализация простого HTTP-сервера с in-memory хранилищем задач.

## Структура проекта

```
pz3-http/
├─ cmd/server/main.go           # точка входа, маршрутизация, graceful shutdown
├─ internal/api/
│  ├─ handlers.go              # обработчики HTTP-запросов
│  ├─ handlers_test.go         # юнит-тесты через httptest
│  ├─ middleware.go            # Logging, CORS
│  └─ responses.go             # хелперы JSON/статусы
├─ internal/storage/
│  └─ memory.go                # in-memory хранилище + ErrNotFound
└─ go.mod
```

## Установка и запуск

Требования: Go 1.22+ (в проекте используется методический `go 1.27.1`).

Клонировать/перейти в корень проекта:

```bash
cd pz3-http
```

Запуск:

```bash
go run ./cmd/server
```

Сборка бинарника:

```bash
go build -o bin/server ./cmd/server
```

## API

| Метод | Путь | Описание |
|---|---|---|
| GET | `/health` | Проверка работы сервера |
| GET | `/tasks` | Список задач. Поддержка фильтра `?q=подстрока` |
| POST | `/tasks` | Создание задачи. Тело `{"title": "..."}`. Принимает `Content-Type: application/json` |
| GET | `/tasks/{id}` | Получение задачи по ID |
| PATCH | `/tasks/{id}` | Отметка выполненной. Текущая реализация всегда переводит в `Done=true` |
| DELETE | `/tasks/{id}` | Удаление задачи (возвращает `204 No Content`) |

Примеры тестирования (curl):

```bash
# health
curl -i http://localhost:8080/health

# создать задачу
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Написать отчет"}'

# список
curl -i "http://localhost:8080/tasks?q=отчет"

# получить по id
curl -i http://localhost:8080/tasks/1

# отметить выполненной
curl -i -X PATCH http://localhost:8080/tasks/1

# удалить
curl -i -X DELETE http://localhost:8080/tasks/1
```

## Валидация

- `title` нельзя пустым (после `TrimSpace`) — `400 Bad Request`
- длина `title` в символах (руны `[]rune`) не больше 140 — `422 Unprocessable Entity`
- неверный `id` — `400 Bad Request`
- задача не найдена — `404 Not Found`
- неверный JSON — `400 Bad Request`

## Особенности

- **In-memory**: данные теряются при перезапуске сервера.
- **Паттерн `statusRecorder`** в middleware: позволяет логировать статус-код ответа (не только факт вызова).
- **Graceful shutdown**: по `Ctrl+C`/`SIGTERM` сервер перестаёт принимать новые подключения и даёт текущим запросам завершиться в течение 5 секунд.
- **CORS**: минимальные заголовки (`Access-Control-Allow-Origin: *`, `Access-Control-Allow-Methods: GET,POST,PATCH,DELETE,OPTIONS`, `Access-Control-Allow-Headers: Content-Type`). Preflight `OPTIONS` отвечает `204` без проксирования.
- **Path-парсинг вручную**: для `/tasks/{id}` без сторонних роутеров (используется `strings.Split`).

## Тесты

```bash
go test -v ./internal/api
```

Пакет `internal/api` содержит юнит-тесты через `httptest` (Create OK/422, Get 404, Patch Done, Delete 204).