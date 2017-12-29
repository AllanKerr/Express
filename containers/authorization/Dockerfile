FROM golang:alpine

RUN apk add --update git

ADD ./app "$GOPATH/src/app"
WORKDIR "$GOPATH/src/app"
RUN go get
RUN go install app

EXPOSE 8080
