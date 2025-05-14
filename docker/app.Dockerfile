FROM golang:1.24

WORKDIR /app

# ⬅️ move this up so it always busts the cache
COPY . .

RUN go mod download

# Install migrate binary (optional dev tool)
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz \
  | tar xvz && mv migrate /usr/local/bin/migrate

RUN apt-get update && apt-get install -y bash

RUN go build -o server ./cmd/server

# Install golangci-lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

CMD ["./server"]
