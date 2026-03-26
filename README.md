# Email Sender Service

Go service that consumes emails from a RabbitMQ queue and sends them via SMTP.

## Structure

```
emailsender/
├── cmd/emailsender/main.go          — Entrypoint
├── internal/
│   ├── config/config.go             — Configuration from .env
│   ├── consumer/consumer.go         — RabbitMQ Consumer with auto-reconnect
│   ├── mailer/mailer.go             — SMTP email sending
│   └── template/
│       ├── template.go              — HTML template rendering (embed.FS)
│       └── email.html               — Email template
├── go.mod
└── .env
```

## Setup

```bash
cd emailsender
go mod tidy
cp .env.example .env   # edit .env as needed
go run cmd/emailsender/main.go
```

## Environment Variables (.env)

| Variable | Description |
|----------|-------------|
| `RABBITMQ_URL` | AMQP connection string |
| `SMTP_HOST` | SMTP server host |
| `SMTP_PORT` | SMTP server port |
| `SMTP_USER` | SMTP username |
| `SMTP_PASS` | SMTP password |
| `SMTP_FROM` | Sender email address |
| `SMTP_FROM_NAME` | Sender display name |

## Message Format (RabbitMQ Queue: `mail`)

```json
{
  "to": "user@example.com",
  "subject": "Код входа",
  "code": "123456",
  "magicUrl": "https://example.com/login-email/magic/abc123",
  "websiteName": "Avtodnr.ru",
  "websiteUrl": "https://avtodnr.ru"
}
```

## Features

- Auto-reconnect on RabbitMQ connection loss (5s backoff)
- Graceful shutdown (SIGINT/SIGTERM)
- Prefetch 1 (one message at a time)
- Ack after successful send, Nack+Requeue on failure
- HTML templates via Go embed.FS

---

# Email Sender Service (DE)

Go-Service der Emails aus einer RabbitMQ-Queue konsumiert und via SMTP versendet.

## Struktur

```
emailsender/
├── cmd/emailsender/main.go          — Entrypoint
├── internal/
│   ├── config/config.go             — Konfiguration aus .env
│   ├── consumer/consumer.go         — RabbitMQ Consumer mit Auto-Reconnect
│   ├── mailer/mailer.go             — SMTP Email-Versand
│   └── template/
│       ├── template.go              — HTML-Template Rendering (embed.FS)
│       └── email.html               — Email-Template
├── go.mod
└── .env
```

## Setup

```bash
cd emailsender
go mod tidy
cp .env.example .env   # oder .env anpassen
go run cmd/emailsender/main.go
```

## Umgebungsvariablen (.env)

| Variable | Beschreibung |
|----------|-------------|
| `RABBITMQ_URL` | AMQP Connection-String |
| `SMTP_HOST` | SMTP Server Host |
| `SMTP_PORT` | SMTP Server Port |
| `SMTP_USER` | SMTP Username |
| `SMTP_PASS` | SMTP Passwort |
| `SMTP_FROM` | Absender Email-Adresse |
| `SMTP_FROM_NAME` | Absender Name |

## Message-Format (RabbitMQ Queue: `mail`)

```json
{
  "to": "user@example.com",
  "subject": "Код входа",
  "code": "123456",
  "magicUrl": "https://example.com/login-email/magic/abc123",
  "websiteName": "Avtodnr.ru",
  "websiteUrl": "https://avtodnr.ru"
}
```

## Features

- Auto-Reconnect bei RabbitMQ-Verbindungsverlust (5s Backoff)
- Graceful Shutdown (SIGINT/SIGTERM)
- Prefetch 1 (eine Message gleichzeitig)
- Ack nach erfolgreichem Senden, Nack+Requeue bei Fehler
- HTML-Templates via Go embed.FS

---

# Сервис отправки Email (RU)

Go-сервис, который потребляет сообщения из очереди RabbitMQ и отправляет их по SMTP.

## Структура

```
emailsender/
├── cmd/emailsender/main.go          — Точка входа
├── internal/
│   ├── config/config.go             — Конфигурация из .env
│   ├── consumer/consumer.go         — RabbitMQ Consumer с авто-переподключением
│   ├── mailer/mailer.go             — Отправка Email через SMTP
│   └── template/
│       ├── template.go              — Рендеринг HTML-шаблонов (embed.FS)
│       └── email.html               — Шаблон письма
├── go.mod
└── .env
```

## Установка

```bash
cd emailsender
go mod tidy
cp .env.example .env   # отредактируйте .env
go run cmd/emailsender/main.go
```

## Переменные окружения (.env)

| Переменная | Описание |
|----------|-------------|
| `RABBITMQ_URL` | AMQP строка подключения |
| `SMTP_HOST` | Хост SMTP сервера |
| `SMTP_PORT` | Порт SMTP сервера |
| `SMTP_USER` | Имя пользователя SMTP |
| `SMTP_PASS` | Пароль SMTP |
| `SMTP_FROM` | Email адрес отправителя |
| `SMTP_FROM_NAME` | Имя отправителя |

## Формат сообщения (очередь RabbitMQ: `mail`)

```json
{
  "to": "user@example.com",
  "subject": "Код входа",
  "code": "123456",
  "magicUrl": "https://example.com/login-email/magic/abc123",
  "websiteName": "Avtodnr.ru",
  "websiteUrl": "https://avtodnr.ru"
}
```

## Возможности

- Авто-переподключение при потере связи с RabbitMQ (5с задержка)
- Graceful shutdown (SIGINT/SIGTERM)
- Prefetch 1 (одно сообщение за раз)
- Ack после успешной отправки, Nack+Requeue при ошибке
- HTML-шаблоны через Go embed.FS
