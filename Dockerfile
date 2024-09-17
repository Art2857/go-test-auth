FROM golang:1.23.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8000

CMD ["go", "run", "./main.go"]
