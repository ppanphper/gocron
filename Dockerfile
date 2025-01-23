FROM golang:1.23-alpine AS builder

RUN apk update \
    && apk add --no-cache git ca-certificates make bash yarn nodejs

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct



WORKDIR /app

# 复制当前目录下的所有文件到 /app 目录
COPY . /app

RUN make install-vue \
    && make build-vue \
    && make statik \
    && CGO_ENABLED=0 make gocron


# RUN git clone https://github.com/ppanphper/gocron.git \
#     && cd gocron \
#     && make install-vue \
#     && make build-vue \
#     && make statik \
#     && CGO_ENABLED=0 make gocron

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -g app app \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app

COPY --from=builder /app/bin/gocron .

RUN chown -R app:app ./

EXPOSE 5920

USER app

ENTRYPOINT ["/app/gocron", "web"]
