FROM golang:alpine AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

WORKDIR /build
ADD go.mod .
ADD go.sum .
ADD main.go .
COPY . .
RUN go build -o main

FROM scratch
WORKDIR /app
COPY application.yml .
COPY --from=builder /build/main /app

EXPOSE 8081
CMD ["./main"]
