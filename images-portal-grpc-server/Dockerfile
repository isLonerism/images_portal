FROM golang:latest as builder
WORKDIR $GOPATH/src/github.com/bsuro10/images-portal/images-portal-grpc-server
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /build/server ./server

FROM docker:dind
USER root
COPY --from=builder /build/server /app/
WORKDIR /app
EXPOSE 7777
CMD ["./server"]
