FROM alpine:latest

ENV GOCRON_NODE_VERSION=v2.0.1

RUN apk add --no-cache ca-certificates  tzdata bash \
    &&  mkdir -p /app \
    &&  wget -P /tmp  https://github.com/ppanphper/gocron/releases/download/${GOCRON_NODE_VERSION}/gocron-${GOCRON_NODE_VERSION}-linux-amd64.tar.gz \
    &&  cd /tmp \
    &&  tar  zvxf gocron-${GOCRON_NODE_VERSION}-linux-amd64.tar.gz  \
    &&  mv /tmp/gocron-linux-amd64/gocron /app \
    &&  rm  -rf /tmp/* \
    &&  cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app
EXPOSE 5920

ENTRYPOINT ["/app/gocron"]
