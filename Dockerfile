FROM golang:1.17.2-alpine3.14  AS builder
WORKDIR /apps
COPY . /apps/
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk --no-cache add gcc musl-dev
RUN go env -w GOPROXY=https://goproxy.cn && \
    go mod tidy && \
    go build -o bin/mydis cmd/mydis/*.go
FROM alpine:3.14 AS runner
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk --no-cache add tzdata && \
    rm -rf /var/cache/apk/*
WORKDIR /apps
COPY --from=builder /apps/bin/mydis /apps/bin/mydis
COPY conf /apps/
EXPOSE 6380
ENTRYPOINT [ "/apps/bin/mydis", "-c", "/apps/conf" ]