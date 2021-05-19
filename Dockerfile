FROM golang:1.14.4-alpine3.12 as builder

# Dev: pass this to `docker build`
# --build-arg APK_REPO=mirrors.aliyun.com
#ARG APK_REPO=dl-cdn.alpinelinux.org
#RUN sed -i "s|//dl-cdn.alpinelinux.org|//${APK_REPO}|g" /etc/apk/repositories; \
#    apk add --no-cache make curl git build-base

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o bin/feishu .

FROM scratch

COPY --from=builder /app/bin/feishu /app

ENTRYPOINT [ "/app" ]
