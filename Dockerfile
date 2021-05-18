FROM golang:1.14.4-alpine3.12 as builder

ARG APK_REPO=mirrors.aliyun.com
RUN sed -i "s|//dl-cdn.alpinelinux.org|//${APK_REPO}|g" /etc/apk/repositories; \
    apk add --no-cache make curl git build-base

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build

FROM scratch

WORKDIR /app
COPY --from=builder /app/actions-feishu /app/actions-feishu

ENTRYPOINT [ "/app/actions-feishu" ]
