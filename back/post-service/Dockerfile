FROM golang:1.21-alpine3.19

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/main.go

EXPOSE 8090

CMD ["app"]