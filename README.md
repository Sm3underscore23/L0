# L0
## Описание проекта

---
## Оглавление
- [Технологии](#технологии)
- [Начало работы](#начало-работы)
- [Основные маршруты API](#основные-маршруты-api)
- [Архитектура и реализация](#архитектура-и-реализация)
- [Валидация](#валидация)
- [Тестирование](#тестирование)
---

## Технологии
### Основные
- [Go](https://github.com/golang/go)
- [PostgreSQL](https://github.com/postgres/postgres)
- [Docker](https://github.com/docker)
- [Docker compose](https://github.com/docker/compose)

### Вспомогательные
- [Chi](https://github.com/go-chi/chi)
- [PGX](https://github.com/jackc/pgx)
- [Goose](https://github.com/pressly/goose)
- [Squirrel](https://github.com/Masterminds/squirrel)
- [Sarama](https://github.com/IBM/sarama)
- [LRU cache](https://github.com/hashicorp/golang-lru)
---

## Начало работы

### Конфигурация проекта

Проект использует два источника конфигурации:

#### 1. YAML-файл

```yaml
# Основные настройки HTTP сервера
http:
  host: localhost       # Хост для запуска HTTP сервера
  port: 8080            # Порт для HTTP соединений

# Настройки кэширования данных
cache:
  limit: 10             # Максимальное количество элементов в кэше
  recover_limit: 5      # Лимит для восстановления кэша

# Настройки логирования
logger:
  log_level: debug      # Уровень детализации логов: debug, info, warn, error

# Настройки подключения к PostgreSQL
postgres:
  host: localhost       # Хост базы данных PostgreSQL
  port: 5432            # Порт PostgreSQL
  db_name: taskL0DB     # Название базы данных
  user: admin           # Имя пользователя для подключения
  sslmode: disable      # Режим SSL: disable, require, verify-full и т.д.
  migrations_dir: migrations  # Директория с миграциями базы данных

# Настройки Kafka consumer
kafka:
  consumer_group: my-consumer-group  # Идентификатор группы потребителей
  broker_list:
    - localhost:9092    # Список брокеров Kafka
  topic: order-topic    # Название топика для подписки
  workers_num: 10       # Количество worker'ов для обработки сообщений
```

---
#### 2. Переменные окружения (`.env`)

```env
DB_PASSWORD=1234                     # Пароль для пользователя базы данных
```

### Запуск проекта
Создайте файл _.env_, (образец _.env.example_)
Отредактируйте или создайте (по необходимости) конфиг в папке _./config_

#### PostgreSQL и Kafka:
Запуск:
```bash
 make compose-up
```
Удаление контейнеров:
```bash
make compose-down
```

#### Запускт backend приложения локально:
```bash
make run
```

### Frontend веб-интерфейса:
Для запука frontend веб-интерфейса перейдите в папку interface и запустите:
```bash
npm install
API_BASE_URL=http://your-backend-url node server.js
```
Где _your-backend-url_ - host:port backend части (example: localhost:8080)

#### Test_producer:
Для запука тестового producer для генерации заказов перейдите в папку test_producer и запустите:
```bash
go mod download
go run cmd/main.go -config-path your_config_path -order-num ordernum
```
Где:
_your_config_path_ - путь до конфига с информацией о Kafka (example: config/local_config.yaml)
_ordernum_ - количество заказов для генерации и отправки (example: 10)

## Основные маршруты API
### Backend

##


## Архитектура и реализация

```

```

## Валидация

## Тестирование

### 1. Unit

```bash
make test
```

## 

### Валидация

См. docs/validation.md
