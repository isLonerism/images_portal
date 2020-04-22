FROM golang:1.12.1 as builder
WORKDIR $GOPATH/src/github.com/bsuro10/images_portal/images-portal-grpc-server
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /build/server ./server

FROM docker:1.12.1-dind
USER root
ENV DOCKER_HOST="tcp://localhost:2375" \ 
    DOCKER_API_VERSION="1.24"
ENTRYPOINT ["/bin/sh" ,"-c"]
WORKDIR /app
EXPOSE 7777
CMD ["/usr/local/bin/dockerd-entrypoint.sh & ./server"]
COPY --from=builder /build/server /app/
