# Система обработки ошибок в API

## Концепция

Система обработки ошибок в API построена на основе кодов ошибок, которые возвращаются клиенту. Каждый код ошибки соответствует определенной ситуации и помогает клиенту понять, что именно произошло.

## Структура ответа API

Все ответы API имеют следующую структуру:

```json
{
  "success": true|false,
  "data": {}, // только для успешных ответов
  "error": {  // только для ошибок
    "code": "ERROR_CODE",
    "message": "Сообщение об ошибке"
  }
}
```

## Коды ошибок

Коды ошибок разделены на несколько категорий:

### Общие коды ошибок

- `UNKNOWN_ERROR` - неизвестная ошибка
- `INVALID_REQUEST` - неверный формат запроса
- `INTERNAL_ERROR` - внутренняя ошибка сервера
- `NOT_FOUND` - ресурс не найден
- `UNAUTHORIZED` - неавторизованный доступ
- `FORBIDDEN` - доступ запрещен

### Пользовательские коды ошибок

- `USER_NOT_FOUND` - пользователь не найден
- `USER_ALREADY_EXISTS` - пользователь уже существует
- `INVALID_PASSWORD` - неверный пароль
- `INVALID_EMAIL` - неверный формат email
- `INVALID_USERNAME` - неверный формат имени пользователя

### Коды ошибок для операций с данными

- `DATA_NOT_FOUND` - данные не найдены
- `DATA_INVALID` - недопустимые данные
- `DATA_CONFLICT` - конфликт данных

## HTTP статусы

Каждый код ошибки соответствует определенному HTTP-статусу:

- `UNKNOWN_ERROR` - 500 Internal Server Error
- `INVALID_REQUEST` - 400 Bad Request
- `NOT_FOUND` - 404 Not Found
- `UNAUTHORIZED` - 401 Unauthorized
- `FORBIDDEN` - 403 Forbidden
- `USER_NOT_FOUND` - 404 Not Found
- `USER_ALREADY_EXISTS` - 409 Conflict
- `INVALID_PASSWORD` - 400 Bad Request
- `DATA_NOT_FOUND` - 404 Not Found
- `DATA_INVALID` - 400 Bad Request
- `DATA_CONFLICT` - 409 Conflict

## Использование в коде

### Создание новой ошибки

```go
// Создание новой ошибки
err := customerrors.New(customerrors.CodeUserNotFound, "Пользователь не найден")

// Создание ошибки на основе существующей
err := customerrors.NewWithError(originalError, customerrors.CodeInternalError, "Внутренняя ошибка сервера")
```

### Отправка ответа с ошибкой

```go
// Отправка ответа с ошибкой
customerrors.RespondWithError(c, err)

// Отправка успешного ответа
customerrors.RespondWithSuccess(c, data)
```

### Проверка типа ошибки

```go
// Проверка типа ошибки
if customerrors.IsErrorCode(err, customerrors.CodeUserNotFound) {
    // Обработка ошибки "Пользователь не найден"
}
```

## Swagger документация

API документация доступна через Swagger UI. После запуска сервера, документация будет доступна по адресу:

```
http://localhost:8080/swagger/index.html
```

Swagger предоставляет:
- Интерактивную документацию всех API endpoints
- Возможность тестирования API прямо из браузера
- Описание всех моделей данных и параметров запросов
- Информацию о кодах ошибок и их значениях

## Мониторинг и метрики

Система включает в себя встроенную поддержку метрик и мониторинга через Grafana.

### Grafana Dashboard

После запуска инфраструктуры через Docker Compose, Grafana будет доступна по адресу:

```
http://localhost:3000
```

**Данные для входа по умолчанию:**
- Логин: `admin`
- Пароль: `admin`

### Источник метрик

Grafana получает метрики напрямую от приложения через endpoint `/metrics` (порт 8081). Приложение экспортирует метрики в формате Prometheus.

### Переменные окружения для Grafana

Для настройки Grafana используются следующие переменные окружения:

```bash
# Учетные данные администратора Grafana
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=admin
```

### Доступные метрики

Система собирает следующие метрики (namespace "via"):

- **via_requests_per_second** - текущий RPS по endpoint'ам и HTTP методам
- **via_request_duration_seconds** - гистограмма времени выполнения запросов
- **via_request_total** - общее количество запросов
- **via_error_total** - счетчики ошибок по типам
- **via_circuitbreaker_total** - метрики Circuit Breaker
- **via_storage_data** - метрики использования хранилища
- **via_out_request_total** - исходящие запросы

### Запуск мониторинга

Для запуска полной инфраструктуры с мониторингом:

```bash
docker compose up -d
```

Grafana автоматически настроится с источником данных от приложения и будет готова к использованию. Рекомендуется изменить пароль администратора при первом входе в систему.

**Подробное руководство по настройке и использованию мониторинга см. в [MONITORING.md](MONITORING.md)** 