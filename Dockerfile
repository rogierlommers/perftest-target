FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o perftest-target .

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/perftest-target .
COPY --from=builder /app/web/template ./web/template

EXPOSE 8080

ENTRYPOINT ["./perftest-target"]
