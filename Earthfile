VERSION 0.8


go-image:
    FROM golang:1.22-alpine
    RUN apk add --no-cache upx

build:
    FROM +go-image

    WORKDIR /workdir
    
    COPY cmd cmd
    COPY pkg pkg
    COPY go.mod go.mod
    COPY go.sum go.sum
    RUN go mod tidy
    RUN go mod download
    RUN go build -ldflags="-w -s" -o imagectl cmd/main.go
    RUN upx -1 imagectl

    SAVE ARTIFACT imagectl AS LOCAL ./build/imagectl