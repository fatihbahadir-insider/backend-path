FROM golang:1.25.5-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/rubenv/sql-migrate/...@latest

COPY . .

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o /app/bin/app ./main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/bin/app /app/app

COPY --from=builder /app/dbconfig.yml /app/dbconfig.yml

COPY --from=builder /app/database/migrations /app/database/migrations

COPY --from=builder /go/bin/sql-migrate /home/nonroot/go/bin/sql-migrate

ENV APP_PORT=8080

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/app/app"]
