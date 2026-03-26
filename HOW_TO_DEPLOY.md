# How to Deploy Email Sender

## Build

```bash
CGO_ENABLED=0 go build -o emailsender ./cmd/emailsender/
```

## Configuration

```bash
cp .env.example .env
# Edit .env with your SMTP and RabbitMQ credentials
```

Required: `RABBITMQ_URL`
Optional: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`, `SMTP_FROM_NAME`

The service loads `.env` automatically, or falls back to environment variables.

## Option 1: Systemd (recommended for single server)

1. Copy the binary:

```bash
sudo cp emailsender /usr/local/bin/emailsender
sudo mkdir -p /opt/emailsender
sudo cp .env /opt/emailsender/.env
```

2. Create a dedicated user:

```bash
sudo useradd -r -s /usr/sbin/nologin emailsender
sudo chown emailsender:emailsender /opt/emailsender/.env
sudo chmod 600 /opt/emailsender/.env
```

3. Create the systemd unit at `/etc/systemd/system/emailsender.service`:

```ini
[Unit]
Description=Email Sender Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/emailsender
WorkingDirectory=/opt/emailsender
EnvironmentFile=/opt/emailsender/.env
Restart=always
RestartSec=5
User=emailsender

[Install]
WantedBy=multi-user.target
```

4. Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now emailsender
```

5. Check status and logs:

```bash
sudo systemctl status emailsender
journalctl -u emailsender -f
```

## Option 2: Docker

1. Create a `Dockerfile`:

```dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o emailsender ./cmd/emailsender/

FROM alpine:3.20
COPY --from=build /app/emailsender /emailsender
CMD ["/emailsender"]
```

2. Build and run:

```bash
docker build -t emailsender .
docker run -d --restart=always --env-file .env --name emailsender emailsender
```

3. Check logs:

```bash
docker logs -f emailsender
```

## Notes

- The service has built-in RabbitMQ auto-reconnect (5s delay between retries)
- Graceful shutdown on SIGINT/SIGTERM
- Prefetch is set to 1, so no messages are lost on shutdown
- Failed emails are requeued automatically via NACK
