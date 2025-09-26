FROM golang:1.21.4-alpine

WORKDIR /

RUN apk add build-base python3-dev py3-pip git squashfs-tools
RUN git clone https://github.com/ReFirmLabs/binwalk.git && \
    cd binwalk && \
    git checkout cddfede795971045d99422bd7a9676c8803ec5ee && \
    python3 setup.py install && \
    cd /
    
COPY ./config.yaml .
COPY ./services/cmd/core/ ./cmd/core/
COPY ./services/common/ ./common/
COPY ./services/config/ ./config/

COPY ./services/core/ ./core/
COPY ./services/pkg/ ./pkg/
COPY ./services/utils/ ./utils/
RUN rm -r ./pkg/packets

COPY ./services/go.mod ./services/go.sum ./

RUN go mod download && go mod verify

RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/core/main.go
RUN mkdir files

ENTRYPOINT [ "./main" ]