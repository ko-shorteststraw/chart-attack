FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /chart-attack ./cmd/server
RUN go build -o /chart-attack-seed ./cmd/seed

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /chart-attack /app/chart-attack
COPY --from=builder /chart-attack-seed /app/chart-attack-seed
COPY migrations/ /app/migrations/
COPY templates/ /app/templates/
COPY static/ /app/static/

EXPOSE 8080

CMD sh -c "/app/chart-attack-seed && /app/chart-attack"
