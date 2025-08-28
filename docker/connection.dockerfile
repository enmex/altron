FROM golang:1.21.4-alpine
WORKDIR /

RUN apk add gcc libc-dev

COPY ./config.yaml .
COPY ./services/cmd/connection/ ./cmd/connection/
COPY ./services/common ./common/
COPY ./services/config/app_config.go ./config/
COPY ./services/config/amqp_config.go ./config/
COPY ./services/config/redis_config.go ./config/
COPY ./services/connection ./connection/
COPY ./services/pkg ./pkg/
COPY ./services/utils ./utils/
RUN rm -r ./pkg/packets

COPY ./services/go.mod ./services/go.sum ./

RUN go mod download && go mod verify

RUN go build -o main ./cmd/connection/main.go

ENTRYPOINT [ "./main" ]
