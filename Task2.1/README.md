# Задание 2.1 — API-тесты микросервиса объявлений

Тестовое задание на стажировку Avito QA (весенняя волна 2026).
Автотесты для REST API [qa-internship.avito.com](https://qa-internship.avito.com), покрывают основные CRUD-сценарии через client/managers/tests (слоистую) архитектуру.

## Стек

- Go 1.25+
- testify
- ozontech/allure-go
- apitest (steinfletcher)
- golangci-lint v2
- goimports

## Архитектура

Решение построено на трёхслойной архитектуре:

- **client** — низкоуровневые HTTP-обёртки над `apitest`, возвращают `*http.Response`
- **managers** — бизнес-логика поверх клиентов, возвращают строки/модели, принимают `testing.TB`
- **tests** — сценарии, работают только через менеджеры и assertion helpers

## Структура проекта

```
Task2.1/
├── internal/
│   ├── client/http/items/      # HTTP-клиент — обёртка над apitest
│   │   └── client.go
│   ├── helpers/
│   │   ├── api-runner/         # Сборка HTTP-клиента (host, debug)
│   │   │   └── runner.go
│   │   └── http-helper/        # AssertStatusCode, ReadResponseBody
│   │       └── response_helper.go
│   ├── managers/items/         # Бизнес-логика поверх клиента
│   │   ├── manager_items.go
│   │   └── models/             # request/response structs
│   │       ├── request.go
│   │       └── response.go
│   └── utils/                  # env, logger, random sellerID
│       ├── env.go
│       ├── logger.go
│       └── random.go
├── tests/scenarios/items/      # Тест-сценарии (allure suite)
│   └── items_test.go
├── .golangci.yml               # Конфиг линтера
├── go.mod
├── go.sum
├── TESTCASES.md
├── BUGS.md
└── allure-example.html         # Пример Allure-отчёта
```

## Переменные окружения

По умолчанию тесты подключаются к `https://qa-internship.avito.com`. Значение берётся из файла `.env` в корне модуля:

```env
API_URL=https://qa-internship.avito.com
```

Для локального оверрайда создайте `.env.override` рядом — он перекрывает `.env` без изменения основного файла:

```env
API_URL=https://your-local-host
```

## Локальный запуск

```bash
git clone https://github.com/IlyaKarzhov/QA-TA-spring-2026-IlyaKarzhov.git
cd QA-TA-spring-2026-IlyaKarzhov/Task2.1

# Установить зависимости
go mod download

# Запустить тесты
go test -v ./tests/scenarios/items/...
```

> Два теста упадут намеренно — они фиксируют баги сервера (BUG-1, BUG-2). Это ожидаемое поведение.

## Allure-отчёт

После запуска тесты сохраняют результаты в `tests/scenarios/items/allure-results/`.

```bash
# Сгенерировать отчёт
allure generate tests/scenarios/items/allure-results --clean -o allure-report

# Открыть в браузере
allure open allure-report
```

Установка Allure CLI:

```bash
# macOS
brew install allure

# Windows (через Scoop — https://scoop.sh/)
scoop install allure
```

Пример отчёта: `allure-example.html`

## Линтер и форматтер

Конфигурация: `.golangci.yml`
Включены: `errcheck`, `gosec`, `govet`, `ineffassign`, `revive`, `staticcheck`, `unused`, `whitespace`
Форматтер: `goimports`

```bash
# Установить goimports
go install golang.org/x/tools/cmd/goimports@latest

# Форматирование
goimports -w ./internal ./tests

# Линтер
golangci-lint run ./...
```

## Покрытие

**Создание объявления (`POST /api/1/item`)**
- Создание с валидными данными
- Граничные значения `price`, `name`, `sellerID`
- Спецсимволы в `name`
- Идемпотентность: два одинаковых запроса создают разные объявления
- Отсутствие обязательных полей (`name`, `price`, `sellerID`)
- Нулевая и отрицательная цена (фиксируют BUG-1 и BUG-2)

**Получение объявления (`GET /api/1/item/:id`)**
- Получение по существующему ID
- Несуществующий ID ->  404

**Список объявлений продавца (`GET /api/1/:sellerID/item`)**
- Продавец с объявлениями — оба UUID присутствуют в ответе
- Продавец без объявлений ->  пустой массив

**Статистика (`GET /api/1/statistic/:id`)**
- Получение по существующему ID
- Несуществующий ID ->  404

**E2E**
- Создание ->  получение по ID ->  получение статистики
