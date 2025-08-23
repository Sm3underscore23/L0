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
/v1/order/find_order/{order_id} - 

/v1/test/get_cache - 

### Frontend
/find-order - 

/order/{order_id} - 

---
## Архитектура и реализация

Проект построен по принципам **чистой архитектуры**, в которой слои логически и физически отделены друг от друга. Каждый уровень зависит только от уровней ниже, что обеспечивает тестируемость, масштабируемость и читаемость.

Основные принципы:
* **Модели и бизнес-логика** изолированы от фреймворков и библиотек.
* **Внешние интерфейсы (API, базы данных)** зависят от внутренней логики, но не наоборот.
* **Связи направлены внутрь** — только внутрь слоя (Dependency Rule).

#### Директори и файлы

**`cmd/app/main.go`** - Точка входа основного приложения
- Парсинг флагов командной строки
- Загрузка переменных окружения
- Инициализация конфигурации

**`config/`** - Конфигурационные файлы**
- `local_config.yaml` - для локальной разработки
- `docker_config.yaml` - для docker-окружения

**`internal/`** - Внутренние пакеты основного приложения**

**`internal/app/app.go`** - Инициализация зависимостей (Dependency Injection)
- Создание логгера, подключение к БД, Kafka, кэшу
- Инициализация репозиториев, use cases, контроллеров
- Запуск HTTP сервера и потребителей Kafka
- Graceful shutdown

**`internal/config/`** - Конфигурация приложения
- `config.go` - структуры и загрузка конфигурации

**`internal/controller/http/`** - HTTP контроллеры**
- `v1/` - API v1 endpoints
  - `get_order_info.go` - получение информации о заказе
  - `get_cache.go` - получение кэшированных данных
  - `middleware/logger.go` - middleware для логирования
- `router.go` - маршрутизация запросов

**`internal/entity/`** - Бизнес-сущности
- `order/` - структура заказа и валидация
- `custom_errors/` - кастомные ошибки приложения
- `logger_tag/` - тэги для логов

**`internal/repository/`** - Слой работы с данными**
- `contracts.go` - интерфейсы репозиториев
- `order/` - реализация репозитория заказов
  - `create.go` - создание заказа
  - `get.go` - получение заказа по ID
  - `get_lasts.go` - получение последних заказов
  - `is_exists.go` - проверка существования заказа

**`internal/usecase/`** - Бизнес-логика**
- `contracts.go` - интерфейсы use cases
- `order/` - реализация бизнес-логики заказов
  - `get_info.go` - получение информации о заказе
  - `get_cache.go` - работа с кэшем
  - `handle_order.go` - обработка входящих заказов
- `cache-loader/` - прогрев кэша

**`pkg/`** - Переиспользуемые пакеты**

**`pkg/httpserver/`** - HTTP сервер с graceful shutdown**
- `server.go` - реализация сервера
- `options.go` - конфигурационные опции

**`pkg/kafka/consumer/`** - Kafka consumer с worker pool**
- `consumer.go` - основной потребитель
- `handler.go` - обработчики сообщений
- `wp.go` - пул воркеров
- `options.go` - конфигурация

**`pkg/postgres/`** - PostgreSQL клиент с миграциями**
- `postgres.go` - подключение к БД
- `migrate.go` - применение миграций
- `options.go` - конфигурация

**`pkg/logger/`** - Логгер (обертка для zap/slog)**

**`interface/`** - Веб-интерфейс на Node.js**
- `server.js` - Express.js сервер
- `views/` - EJS шаблоны
- `package.json` - зависимости Node.js

**`test_producer/`** - Тестовый продюсер заказов**
- Генерация тестовых заказов в Kafka
- Отдельный Go модуль

**`migrations/`** - Миграции базы данных**
- `001_orders.sql` - таблица заказов
- `002_delivery.sql` - данные доставки
- `003_payments.sql` - платежная информация
- `004_items.sql` - товары в заказе

**`testdata/`** - Тестовые данные**
- Конфигурации для тестирования

### Docker
**`Dockerfile`** - Сборка Go приложения**

**`docker-compose.yaml`** - Основной compose файл**
- PostgreSQL, Kafka, Zookeeper

**`app.docker-compose.yaml`** - Запуск приложения**
- Основное приложение, PostgreSQL, Kafka, Zookeeper

## Валидация



## Тестирование

### 1. Unit

```bash
make test
```

## 

### Валидация

См. docs/validation.md
