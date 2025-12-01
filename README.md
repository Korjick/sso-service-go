# SSO Service Go

SSO Service Go — это сервис единого входа (SSO), реализованный на языке Go для безопасной аутентификации и авторизации пользователей

## Особенности
- Поддержка JWT-токенов для сессий и аутентификации
- Интеграция с внешними провайдерами (OAuth, OpenID Connect)
- Управление пользователями и ролями
- RESTful API для авторизации
- Кэширование сессий с Redis
- Логирование с Zap

## Архитектура
Проект следует принципам Clean Architecture:

internal/<br>
├── domain/ # Бизнес-логика и сущности<br>
├── application/ # Сервисы<br>
├── infrastructure/ # Репозитории и внешние сервисы<br>
└── presentation/ # HTTP-слой с Gin<br>


## Технологический стек
- **Backend**: Go 1.24+
- **HTTP Framework**: Gin
- **Базы данных**: PostgreSQL (пользователи), Redis (сессии)
- **Аутентификация**: JWT
- **Логирование**: Zap
- **Сборка**: Makefile, Docker

Проект использует Go на 98.7% и Makefile на 1.3%.[attached_file:1]

## Предварительные требования
- Go 1.24+
- Docker и Docker Compose

## Быстрый старт
1. Клонируйте репозиторий:
```shell
git clone https://github.com/Korjick/sso-service-go.git
cd sso-service-go
```

2. Создайте .env файл:
```shell
cp .env.example .env
```

3. Запустите инфраструктуру:
```shell
docker-compose up -d
```

4. Установите зависимости:
```shell
go mod download
```

6. Запустите сервис:
```shell
go run cmd/main.go
```

Сервис доступен на `http://localhost:8080`.

## API Endpoints

### Аутентификация
- `POST /api/login` — вход в систему
- `POST /api/register` — регистрация
- `POST /api/refresh` — обновление токена
- `POST /api/logout` — выход

### OAuth
- `GET /oauth/authorize` — авторизация
- `POST /oauth/token` — получение токена

## Лицензия
MIT License.
