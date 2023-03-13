FROM golang:latest

WORKDIR /build

COPY . .

WORKDIR /build/cmd

RUN go build -o main

EXPOSE 80

ENTRYPOINT [ "./main" ]