# How to Deploy Email Sender

## Configuration

```bash
cp .env.example .env
```

Required: `RABBITMQ_URL`
Optional: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`, `SMTP_FROM_NAME`

## Run

```bash
docker compose up -d --build
```

## Logs

```bash
docker compose logs -f
```

## Restart

```bash
docker compose restart
```

## Stop

```bash
docker compose down
```

## Notes

- The service has built-in RabbitMQ auto-reconnect (5s delay between retries)
- Graceful shutdown on SIGINT/SIGTERM
- Prefetch is set to 1, so no messages are lost on shutdown
- Failed emails are requeued automatically via NACK
