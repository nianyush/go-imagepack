VERSION 0.8


build:
    FROM golang:1.22

    WORKDIR /workdir

    COPY cmd cmd
    COPY pkg pkg
    COPY go.mod go.mod
    COPY go.sum go.sum
    RUN go mod tidy
    RUN go mod download
    RUN go build -o pack cmd/main.go

    SAVE ARTIFACT pack AS LOCAL ./build/pack