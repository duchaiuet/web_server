# syntax=docker/dockerfile:1
FROM golang:1.17-alpine
ADD . /web_server
WORKDIR /web_server

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /account_service

EXPOSE 8088

CMD ["/account_service" ]