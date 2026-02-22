# Org Structure API

REST API для управления организационной структурой компании - подразделениями и сотрудниками.

## Стек

- **Go** (net/http)
- **PostgreSQL**
- **GORM** - ORM
- **goose** - миграции
- **Docker / docker-compose**

## Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone https://github.com/rudolfkova/org-structure-api.git
cd org-structure-api
```

### 2. Создать `.env` файл

```bash
cp .env.example .env
```

Необходимо отредактировать `.env`, задать свои значения для `POSTGRES_USER`, `POSTGRES_PASSWORD` и т.д.

### 3. Запустить

```bash
docker-compose up --build
```

Сервис автоматически:
1. Поднимет PostgreSQL
2. Прогонит миграции через goose
3. Запустит API

## API

### Подразделения

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/departments/` | Создать подразделение |
| `GET` | `/departments/{id}` | Получить подразделение с деревом |
| `PATCH` | `/departments/{id}` | Переименовать / переместить |
| `DELETE` | `/departments/{id}` | Удалить подразделение |

### Сотрудники

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/departments/{id}/employees/` | Создать сотрудника в подразделении |

---

### POST `/departments/`

```json
{
  "name": "Backend",
  "parent_id": 1
}
```

`parent_id` - опционально. Если не указан, создаётся корневое подразделение.

---

### GET `/departments/{id}`

Query-параметры:

| Параметр | Тип | По умолчанию | Описание |
|----------|-----|-------------|----------|
| `depth` | int | `1` | Глубина вложенных подразделений (макс. 5) |
| `include_employees` | bool | `true` | Включать ли сотрудников в ответ |

Пример: `GET /departments/1?depth=3&include_employees=true`

---

### PATCH `/departments/{id}`

```json
{
  "name": "New Name",
  "parent_id": 5
}
```

Оба поля опциональны. Нельзя создать цикл в дереве - вернёт `409 Conflict`.

---

### DELETE `/departments/{id}`

Query-параметры:

| Параметр | Значения | Описание |
|----------|----------|----------|
| `mode` | `cascade` / `reassign` | Режим удаления |
| `reassign_to_department_id` | int | Обязателен при `mode=reassign` |

- `cascade` - удаляет подразделение, всех сотрудников и дочерние подразделения
- `reassign` - переводит всех сотрудников в указанное подразделение, затем удаляет

Пример: `DELETE /departments/3?mode=reassign&reassign_to_department_id=1`

---

### POST `/departments/{id}/employees/`

```json
{
  "full_name": "Иван Петров",
  "position": "Senior Developer",
  "hired_at": "2022-03-15T00:00:00Z"
}
```

`hired_at` - опционально.

## Переменные окружения

| Переменная | Описание | Пример |
|------------|----------|--------|
| `DATABASE_URL` | DSN для подключения к PostgreSQL | `host=db port=5432 user=postgres password=postgres dbname=orgstructure sslmode=disable` |
| `BIND_ADDR` | Адрес и порт сервера | `:8080` |
| `LOG_LEVEL` | Уровень логирования (`debug`, `info`, `warn`) | `info` |
| `POSTGRES_USER` | Пользователь PostgreSQL | `postgres` |
| `POSTGRES_PASSWORD` | Пароль PostgreSQL | `postgres` |
| `POSTGRES_DB` | Имя базы данных | `orgstructure` |

## Тесты

### Unit-тесты (usecase)

Не требуют базы данных, запускаются сразу:

```bash
go test ./internal/test/usecase/... -v
```

### Интеграционные тесты

Требуют запущенной PostgreSQL. Задай переменную окружения с DSN тестовой БД:

```bash
# Linux / macOS
export TEST_DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=orgstructure_test sslmode=disable"
go test ./internal/test/integration/... -v -timeout 60s

# Windows (PowerShell)
$env:TEST_DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=orgstructure_test sslmode=disable"
go test ./internal/test/integration/... -v -timeout 60s
```

Если `TEST_DATABASE_URL` не задан - интеграционные тесты автоматически пропускаются.
