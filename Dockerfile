FROM golang:1.17-alpine3.14  AS builder
WORKDIR /apps
COPY . /apps/
RUN go env -w GOPROXY=https://goproxy.cn && \
    go mod tidy && \
    go build -o bin/mydis cmd/mydis/*.go
FROM alpine:3.14 AS runner
COPY --from=builder /apps/bin/mydis /bin/mydis
EXPOSE 6380
#ENTRYPOINT [ "/bin/mydis", "--driver", "mysql"]