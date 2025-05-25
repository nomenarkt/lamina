FROM golang:1.24.3

WORKDIR /app

# ⬅️ move this up so it always busts the cache
COPY . .

RUN go mod download

# Install migrate binary (optional dev tool)
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz \
  | tar xvz && mv migrate /usr/local/bin/migrate

RUN apt-get update && \
    apt-get install -y --no-install-recommends gnupg ca-certificates && \
    apt-get install -y bash && \
    apt-get clean && rm -rf /var/lib/apt/lists/*


RUN go build -o server ./cmd/server

# Install golangci-lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install air for live reload
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

CMD ["air", "-c", ".air.toml"]

#| Mode        | Command in Dockerfile            | Use Case                        |
#| ----------- | -------------------------------- | ------------------------------- |
#| Production  | `CMD ["./server"]`               | Final binary, no reload         |
#| Development | `CMD ["air", "-c", ".air.toml"]` | Live reload with source mapping |
