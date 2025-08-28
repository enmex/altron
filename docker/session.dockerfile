FROM golang:1.21.4-alpine as build
WORKDIR /

RUN apk add --no-cache build-base libpcap-dev docker openrc

COPY ./services/go.mod ./services/go.sum ./
RUN go mod download && go mod verify

COPY ./config.yaml .
COPY ./services .

RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/session/main.go
RUN mkdir logs pcaps files

ENTRYPOINT [ "./main" ]
