package main

import (
	"fmt"
	"log"
	"net"

	"github.com/bsuro10/images-portal/images-portal-grpc-server/api/docker"
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	s := docker.RegisterDockerService(cli)

	grpcServer := grpc.NewServer()

	docker.RegisterDockerServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
