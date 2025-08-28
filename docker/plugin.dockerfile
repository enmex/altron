FROM golang:1.21.4-alpine
WORKDIR /

ENV PLUGINS_DIRECTORY /plugins

RUN apk add gcc libc-dev

COPY ./config.yaml .
COPY ./services/cmd/plugin/ ./cmd/plugin/
COPY ./services/common ./common/
COPY ./services/config/app_config.go ./config/
COPY ./services/config/plugin_config.go ./config/
COPY ./services/plugin ./plugin/
COPY ./services/pkg ./pkg/
COPY ./services/utils ./utils/
COPY ./plugins/ ./plugins/
RUN rm -r ./pkg/packets

COPY ./services/go.mod ./services/go.sum ./

RUN go mod download && go mod verify

RUN cd plugins && \
    for dir in $(find . -maxdepth 1 -type d ! -path . -exec basename {} \;); \
    do \
        go build \
        -buildmode=plugin -o \
        $(echo "./${dir}/${dir}.so") \
        $(echo "./${dir}/${dir}.go");  \
    done

RUN go build -o main ./cmd/plugin/main.go

ENTRYPOINT [ "./main" ]