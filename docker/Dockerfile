FROM golang:1.18 as builder


RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app
RUN go mod download
RUN CGO_ENABLED=0 go build -o app


FROM alpine:3.15
COPY --from=builder /go/src/app /
ENTRYPOINT [ "/app" ]
CMD [ "run" ]