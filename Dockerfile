FROM golang:1.12-stretch

WORKDIR $GOPATH/src/

COPY . . 

RUN	go get github.com/satori/go.uuid
RUN	go get github.com/gorilla/mux
RUN go get github.com/tidwall/gjson

RUN go build service.go

ARG ACCESS_TOKEN

ENV ACCESS_TOKEN=$ACCESS_TOKEN

EXPOSE 8000

CMD "./service"