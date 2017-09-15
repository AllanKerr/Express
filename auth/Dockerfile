FROM golang:alpine

RUN apk add --update git

ADD ./app "$GOPATH/src/app"

WORKDIR "$GOPATH/src/app"
RUN go get

WORKDIR "$GOPATH"
RUN go build -o main app

EXPOSE 8080
CMD "${GOPATH}/main"
