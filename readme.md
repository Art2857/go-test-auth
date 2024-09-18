
**Установка**

```bash
go mod download
```

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

**Запуск**

Запуск из коробки:
```bash
docker-compose up
```

Запуск для разработки:
```bash
docker-compose -f ./docker-compose.database.yml up -d
go run auth
```

**Запуск тестов**

Предварительная подготовка:
```bash
docker-compose -f ./docker-compose.test-database.yml up -d
```

Запуск e2e тестов:
```bash
go test ./tests
```

**Роутинг**

Роут документации:

http://localhost:8000/docs/

Роут Аутентификации:

http://localhost:8000/token/login

Тело запроса:
```json
{
    "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
}
```

Ответ:
```json
{
   "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTBlZWJjOTktOWMwYi00ZWY4LWJiNmQtNmJiOWJkMzgwYTExIiwiaXAiOiJbOjoxXTo2MzY1MCIsInJlZnJlc2hfdG9rZW5faGFzaCI6IiQyYSQxMCRkdHJVVDJxdlZFOTc0LjFsbXFnQTdPNXhnM0tQL3daM1pGeWFtWW1mMzBJYjFsUWl6bWhYZSIsImV4cCI6MTcyNjYwNTI0OH0.-B0fNb-Ln5N16m_PO-Fea7YfwtjFC6Bo-oq2t04Cq_ME-Xm_skH9YUlEIhOEQoyZSa5SJoZUuSacxuqeFS27QQ",
   "refresh_token": "2tKLjPcGh4j5pXqVpzhSQhF1uQyil5NVEtPMcLQVcPE="
}
```

Роут Обновления пары токенов:

http://localhost:8000/token/refresh

Тело запроса:
```json
{
   "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTBlZWJjOTktOWMwYi00ZWY4LWJiNmQtNmJiOWJkMzgwYTExIiwiaXAiOiJbOjoxXTo2MzY1MCIsInJlZnJlc2hfdG9rZW5faGFzaCI6IiQyYSQxMCRkdHJVVDJxdlZFOTc0LjFsbXFnQTdPNXhnM0tQL3daM1pGeWFtWW1mMzBJYjFsUWl6bWhYZSIsImV4cCI6MTcyNjYwNTI0OH0.-B0fNb-Ln5N16m_PO-Fea7YfwtjFC6Bo-oq2t04Cq_ME-Xm_skH9YUlEIhOEQoyZSa5SJoZUuSacxuqeFS27QQ",
   "refresh_token": "2tKLjPcGh4j5pXqVpzhSQhF1uQyil5NVEtPMcLQVcPE="
}
```

Ответ:
```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTBlZWJjOTktOWMwYi00ZWY4LWJiNmQtNmJiOWJkMzgwYTExIiwiaXAiOiJbOjoxXTo2MzY1MCIsInJlZnJlc2hfdG9rZW5faGFzaCI6IiQyYSQxMCRKRFA0cGNlQUZ1Z3BqeGFHbjJuTzFPSXYvekNWLmRpVHJDcFFIOGszdnM2SEhxSDBHUlpHTyIsImV4cCI6MTcyNjYwNTI2NH0.2INR-GzVLEkgbkIECeherVYcXfZYpnpDgNutU0cFHfF-N5WiciPT5mcQZvCEEI6Hkr1RUE9Njna6A3N1tClRvA",
    "refresh_token": "yDm11GSqlt69qvvHdx1oTVJOXNkUX3RrTAwLx8s3c8M="
}
```

**Техническое задание:**

Написать часть сервиса аутентификации.

Два REST маршрута:

- Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
- Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

**Требования:**

Access токен тип JWT, алгоритм SHA512, хранить в базе строго запрещено.

Refresh токен тип произвольный, формат передачи base64, хранится в базе исключительно в виде bcrypt хеша, должен быть защищен от изменения на стороне клиента и попыток повторного использования.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним.

Payload токенов должен содержать сведения об ip адресе клиента, которому он был выдан. В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные).

**Используемые технологии:**

- Go
- JWT
- PostgreSQL
