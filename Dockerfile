FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o emailsender ./cmd/emailsender/

FROM alpine:3.20
COPY --from=build /app/emailsender /emailsender
CMD ["/emailsender"]
