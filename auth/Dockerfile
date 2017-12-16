FROM golang:alpine

RUN apk add --update git

ADD ./app "$GOPATH/src/app"
ADD ./cli "$GOPATH/src/cli"

WORKDIR "$GOPATH/src/app"
RUN go get
RUN go install app

WORKDIR "$GOPATH/src/cli"
RUN go get
RUN go install cli

EXPOSE 8080
CMD "app"
