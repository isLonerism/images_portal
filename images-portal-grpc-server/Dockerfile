FROM golang:latest as builder
WORKDIR $GOPATH/src/github.com/bsuro10/images-portal/images-portal-grpc-server
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /build/server ./server

FROM docker:dind
USER root
ENTRYPOINT ["/bin/sh" ,"-c"]
WORKDIR /app
EXPOSE 7777
CMD ["/usr/local/bin/dockerd-entrypoint.sh & ./server"]
COPY --from=builder /build/server /app/
