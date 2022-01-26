FROM golang:1.16-alpine

ADD . /app

WORKDIR /app

ENV GO111MODULE=on
ENV GOPATH=""
ENV GOPROXY=off
RUN go build -mod vendor -ldflags "-w -s" -o /tmp/api-gw-backend cmd/main.go

EXPOSE 8080

ENTRYPOINT [ "/tmp/api-gw-backend"]