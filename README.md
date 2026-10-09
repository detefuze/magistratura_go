# myapp

Учебный Go-проект по практическому занятию №2 («Структура Go-проекта»).

Демонстрирует раскладку кода по рекомендованным каталогам: `cmd/` — точка входа,
`internal/` — приватный код приложения, `utils/` — вспомогательные функции.
`main.go` минимален и лишь запускает приложение из `internal/app`.

Сервис поднимает HTTP-сервер на порту `:8080`, логирует все запросы, добавляет к ответам
заголовок `X-Request-Id` и умеет отдавать как обычные JSON-ответы, так и единообразные
JSON-ошибки.

## Возможности

| Маршрут | Метод | Описание                                                                 |
|---------|-------|--------------------------------------------------------------------------|
| `/`     | GET   | Текстовый ответ `Hello, Go project structure!` (`text/plain`)            |
| `/ping` | GET   | JSON `{"status":"ok","time":"<RFC3339>"}` с текущим временем (UTC)       |
| `/fail` | GET   | Пример единого формата ошибки: код `400`, тело `{"error":"bad_request_example"}` |

Дополнительно:

- **`X-Request-Id`** — middleware `withRequestID` берёт заголовок `X-Request-Id` из запроса,
  а если его нет — генерирует 16-символьный hex (`utils.NewID16`). Значение всегда
  возвращается в заголовке ответа, что удобно для трассировки запросов.
- **Логирование** — `utils.LogRequest` пишет в консоль время, адрес клиента, метод и путь;
  `utils.LogInfo` / `utils.LogError` — служебные сообщения с префиксом `[INFO]` / `[ERROR]`.

## Структура проекта

```
myapp/
├─ cmd/
│  └─ myapp/
│     └─ main.go              # точка входа, вызывает app.Run()
├─ internal/
│  └─ app/
│     ├─ app.go               # HTTP-сервер, маршруты, middleware withRequestID
│     └─ handlers/
│        └─ ping.go           # обработчик /ping (вынесен в отдельный пакет)
├─ utils/
│  ├─ logger.go               # LogRequest / LogInfo / LogError / NewID16
│  └─ httpjson.go             # WriteJSON / WriteErr — единый формат ответов и ошибок
├─ .vscode/
│  └─ tasks.json              # задачи VS Code: Run myapp, Build myapp
├─ go.mod
├─ .gitignore
└─ README.md
```

Модуль: `github.com/detefuze/myapp`.

## Запуск

```bash
go run ./cmd/myapp
```

После старта в консоли появится `[INFO] ... Server is starting on :8080`.

## Проверка

Linux / macOS / Git Bash:

```bash
curl -i http://localhost:8080/
curl -i http://localhost:8080/ping
curl -i -H "X-Request-Id: demo-123" http://localhost:8080/ping
curl -i http://localhost:8080/fail
```

Windows PowerShell используйте `curl.exe` (просто `curl` — это алиас на `Invoke-WebRequest`):

```powershell
curl.exe -i http://localhost:8080/
curl.exe -i http://localhost:8080/ping
curl.exe -i -H "X-Request-Id: demo-123" http://localhost:8080/ping
curl.exe -i http://localhost:8080/fail
```

Либо средствами PowerShell:

```powershell
(Invoke-WebRequest http://localhost:8080/ping).Headers
Invoke-RestMethod http://localhost:8080/fail
```

В ответе на `/ping` будет заголовок `X-Request-Id` (переданное значение `demo-123`
или сгенерированный hex).

## Сборка

Linux / macOS:

```bash
go build -o bin/myapp ./cmd/myapp
./bin/myapp
```

Windows:

```powershell
go build -o bin\myapp.exe .\cmd\myapp
.\bin\myapp.exe
```

## Задачи VS Code

В `.vscode/tasks.json` описаны задачи:

- **Run myapp** — `go run ./cmd/myapp`
- **Build myapp** — `go build -o bin/myapp ./cmd/myapp`
