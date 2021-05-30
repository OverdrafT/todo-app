FROM golang:1.14.6-alpine3.12 as builder

COPY go.mod go.sum /go/src/github.com/silverspase/todo/
WORKDIR /go/src/github.com/silverspase/todo/
RUN go mod download
COPY . /go/src/github.com/silverspase/todo/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/todo /go/src/github.com/silverspase/todo/


FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/silverspase/todo/build/todo /usr/bin/todo

ENTRYPOINT ["/usr/bin/todo"]
