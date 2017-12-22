FROM golang:alpine

RUN apk add --update git

ADD ./gateway-controller "$GOPATH/src/gateway-controller"
WORKDIR "$GOPATH/src/gateway-controller"
RUN go get
RUN go install gateway-controller

EXPOSE 8080
