FROM golang:1.14.4-alpine3.12 as builder

ONBUILD ARG CI
ARG APK_REPO=mirrors.aliyun.com
RUN if [ ! $CI ];then \
        sed -i "s|//dl-cdn.alpinelinux.org|//${APK_REPO}|g" /etc/apk/repositories \
    fi \
    apk add --no-cache make curl git build-base

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o bin/feishu .

FROM scratch

COPY --from=builder /app/bin/feishu /app

ENTRYPOINT [ "/app" ]
