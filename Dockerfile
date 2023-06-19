FROM golang:latest

WORKDIR /build

COPY . .

WORKDIR /build/cmd

RUN go build -o main

ENV JWT_SECRET_KEY "asdasasdasasdasasdasasdasaa"
ENV CRYPTER_SECRET_KEY "this is secret key enough 32 bit"
ENV DB_HOST "localhost"
ENV DB_PORT "3306"
ENV DB_USERNAME "root"
ENV DB_PASSWORD "password"
ENV DB_NAME "db"
ENV LOG_LEVEL "info"
ENV VAULT_PASSWORD "VeryHard987*"

EXPOSE 80

ENTRYPOINT [ "./main" ]