FROM golang:1.14.4-alpine3.12 as builder


WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o bin/feishu .

FROM scratch

COPY --from=builder /app/bin/feishu /app

ENTRYPOINT [ "/app" ]
