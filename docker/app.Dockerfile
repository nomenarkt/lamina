FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install migrate binary to assist during dev if needed
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz \
  | tar xvz && mv migrate /usr/local/bin/migrate

RUN apt-get update && apt-get install -y bash
RUN go build -o server ./cmd/server

CMD ["./server"]
