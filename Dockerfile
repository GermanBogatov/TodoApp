FROM golang:alpine as builder

WORKDIR /usr/local/go/src/

ADD app/ /usr/local/go/src/

RUN go clean --modcache
#RUN go mod download
#ADD go.mod .
#ADD go.sum .
#RUN go mod tidy
#RUN go mod vendor
RUN go build -o app cmd/main/app.go

FROM alpine

COPY --from=builder /usr/local/go/src//app /
COPY --from=builder /usr/local/go/src/config.yml /
CMD ["/app"]

#FROM golang:1.18

#WORKDIR go/src/app
#COPY . .

#RUN go get -d -v ./...
#RUN go install -v ./...

#CMD ["app"]






