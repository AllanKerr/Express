FROM golang:alpine

RUN apk add --update git

ADD ./express-controller "$GOPATH/src/express-controller"
WORKDIR "$GOPATH/src/express-controller"
RUN go get
RUN go install express-controller

EXPOSE 8080
