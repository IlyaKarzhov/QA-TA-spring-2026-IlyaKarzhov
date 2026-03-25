# Тест-кейсы: API сервиса объявлений

**Базовый URL:** `https://qa-internship.avito.com`

**Техники тест-дизайна:** классы эквивалентности, анализ граничных значений, негативное тестирование

---

| ID | Заголовок | Предусловия | Шаги | Ожидаемый результат | Статус | Автотест |
|----|-----------|-------------|------|---------------------|--------|----------|
| TC-POS-01 | Создание объявления с валидными данными | — | 1. `POST /api/1/item` с телом `{"name":"Test","price":100,"sellerID":123456,"statistics":{"likes":1,"viewCount":1,"contacts":1}}`<br>2. Проверить HTTP-статус<br>3. Проверить тело ответа<br>4. Проверить заголовок `Content-Type` | 200 OK; тело `{"status":"Сохранили объявление - <uuid>"}`, Content-Type: application/json | Passed | `TestCreateItemPositive` |
| TC-POS-02 | Получение объявления по ID | Создано объявление, известен UUID | 1. `GET /api/1/item/{uuid}`<br>2. Проверить структуру ответа<br>3. Проверить `[0].name`, `[0].price`, `[0].sellerId` | 200 OK; массив из 1 элемента; поля совпадают с данными при создании | Passed | `TestGetItemByIDPositive` |
| TC-POS-03 | Получение всех объявлений продавца | — | 1. `POST /api/1/item` с `sellerID=X` — запомнить UUID-1<br>2. `POST /api/1/item` с `sellerID=X` — запомнить UUID-2<br>3. `GET /api/1/{X}/item`<br>4. Проверить, что UUID-1 и UUID-2 присутствуют в ответе | 200 OK; массив содержит оба созданных объявления | Passed | `TestGetItemsBySellerIDPositive` |
| TC-POS-04 | Получение статистики по объявлению | Создано объявление, известен UUID | 1. `GET /api/1/statistic/{uuid}`<br>2. Проверить структуру ответа<br>3. Проверить поля `[0].viewCount`, `[0].contacts`, `[0].likes` | 200 OK; массив из 1 элемента; все поля числовые и ≥ 0 | Passed | `TestGetStatisticPositive` |
| TC-POS-05 | E2E: создание -> получение -> статистика | — | 1. `POST /api/1/item` — получить UUID<br>2. `GET /api/1/item/{uuid}` — проверить `[0].name`, `[0].price`<br>3. `GET /api/1/statistic/{uuid}` — проверить наличие полей статистики | Все три запроса 200 OK; данные объявления совпадают; статистика возвращается для того же UUID | Passed | `TestE2ECreateGetStatistic` |
| TC-NEG-01 | Создание объявления с пустым `name` | — | 1. `POST /api/1/item` с `"name": ""`<br>2. Проверить HTTP-статус и тело ответа | 400 Bad Request; тело содержит сообщение об ошибке | Passed | `TestCreateItemWithoutName` |
| TC-NEG-02 | Создание объявления без поля `price` | — | 1. `POST /api/1/item` без ключа `price` в теле<br>2. Проверить HTTP-статус и тело ответа | 400 Bad Request; тело содержит сообщение об ошибке | Passed | `TestCreateItemWithoutPriceField` |
| TC-NEG-03 | Создание объявления с `sellerID=0` | — | 1. `POST /api/1/item` с `"sellerID": 0`<br>2. Проверить HTTP-статус и тело ответа | 400 Bad Request; тело содержит сообщение об ошибке | Passed | `TestCreateItemWithoutSellerID` |
| TC-NEG-04 | Создание объявления с `price=0` | — | 1. `POST /api/1/item` с `"price": 0`<br>2. Проверить HTTP-статус | 200 OK; ноль — допустимое значение цены | Failed (BUG-1) | `TestCreateItemWithZeroPrice` |
| TC-NEG-05 | Создание объявления с отрицательным `price` | — | 1. `POST /api/1/item` с `"price": -100`<br>2. Проверить HTTP-статус | 400 Bad Request; отрицательная цена недопустима | Failed (BUG-2) | `TestCreateItemWithNegativePrice` |
| TC-NEG-06 | Создание объявления с `sellerID` вне рекомендованного диапазона | — | 1. `POST /api/1/item` с `"sellerID": 100000`<br>2. Проверить HTTP-статус | 200 OK; диапазон 111111–999999 является рекомендацией, не жёстким ограничением | Passed | `TestCreateItemWithSellerIDOutOfRecommendedRange` |
| TC-NEG-07 | Получение объявления по несуществующему ID | — | 1. `GET /api/1/item/{random-uuid}`, UUID заведомо отсутствует в системе<br>2. Проверить HTTP-статус | 404 Not Found | Passed | `TestGetItemByIDNotFound` |
| TC-NEG-08 | Получение объявлений продавца, у которого нет объявлений | — | 1. `GET /api/1/{unused-sellerID}/item`, sellerID без объявлений<br>2. Проверить тело ответа | 200 OK; пустой массив `[]` | Passed | `TestGetItemsBySellerIDNotFound` |
| TC-NEG-09 | Получение статистики по несуществующему ID | — | 1. `GET /api/1/statistic/{random-uuid}`, UUID заведомо отсутствует<br>2. Проверить HTTP-статус | 404 Not Found | Passed | `TestGetStatisticNotFound` |
| TC-COR-01 | Граничные значения `price`: минимальное (1) и максимальное (INT32_MAX) | — | 1. `POST /api/1/item` с `"price": 1`<br>2. `POST /api/1/item` с `"price": 2147483647`<br>3. Проверить HTTP-статус обоих запросов | 200 OK для обоих значений | Passed | `TestCreateItemWithBoundaryPrice` |
| TC-COR-02 | Граничные значения длины `name`: 1 символ и 255 символов | — | 1. `POST /api/1/item` с `"name": "A"`<br>2. `POST /api/1/item` с `name` из 255 символов<br>3. Проверить HTTP-статус обоих запросов | 200 OK для обоих значений | Passed | `TestCreateItemWithBoundaryName` |
| TC-COR-03 | Граничные значения `sellerID`: 111111 и 999999 | — | 1. `POST /api/1/item` с `"sellerID": 111111`<br>2. `POST /api/1/item` с `"sellerID": 999999`<br>3. Проверить HTTP-статус обоих запросов | 200 OK для обоих значений | Passed | `TestCreateItemWithBoundarySellerID` |
| TC-COR-04 | Идемпотентность: два одинаковых POST-запроса создают два разных объявления | — | 1. `POST /api/1/item` с фиксированными данными — сохранить UUID-1<br>2. Повторить тот же запрос — сохранить UUID-2<br>3. Сравнить UUID-1 и UUID-2 | UUID-1 ≠ UUID-2; каждый вызов создаёт новый ресурс | Passed | `TestIdempotentCreate` |
| TC-COR-05 | Спецсимволы в поле `name` | — | 1. `POST /api/1/item` с `"name": "!@#$%^&*()"`<br>2. `GET /api/1/item/{uuid}`<br>3. Проверить `[0].name` | 200 OK; `[0].name` возвращается без изменений | Passed | `TestCreateItemWithSpecialCharsName` |
