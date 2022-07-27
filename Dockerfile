FROM golang:1.18-alpine3.14 as builder

WORKDIR /usr/local/go/src/

ADD app/ /usr/local/go/src/

ADD go.mod .
ADD go.sum .

RUN go clean --modcache
RUN go mod download
RUN go mod tidy
RUN go build -mod=readonly -o app cmd/main/app.go


FROM alpine

COPY --from=builder /usr/local/go/src/app /
COPY --from=builder /usr/local/go/src/config.yml /
CMD ["app"]







