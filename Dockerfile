FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN mkdir cmd internal config

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./config ./config

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/main.go

ENTRYPOINT ["./main"]