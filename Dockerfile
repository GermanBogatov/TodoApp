FROM golang:alpine as builder

WORKDIR /usr/local/go/src/TodoApp

ADD app/ /usr/local/go/src/TodoApp

RUN go clean --modcache

RUN go mod tidy
RUN go mod vendor
RUN go build -o app cmd/main/app.go

FROM alpine

COPY --from=builder /usr/local/go/src/TodoApp/app /
COPY --from=builder /usr/local/go/src/TodoApp/config.yml /
CMD ["/app"]








