FROM golang:alpine

COPY . /go/src/slave
WORKDIR /go/src/slave
ENTRYPOINT go run main.go
EXPOSE 3000