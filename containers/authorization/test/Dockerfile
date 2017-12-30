FROM golang:alpine

RUN apk add --update git

ADD . "$GOPATH/src/integration"
WORKDIR "$GOPATH/src/integration"
RUN go get
