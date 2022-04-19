FROM golang:1.18-alpine

RUN mkdir -p /go/src/app

WORKDIR /go/src/app

COPY . /go/src/app

RUN go mod download
RUN go build -o app main.go

ENTRYPOINT [ "/go/src/app/app" ]