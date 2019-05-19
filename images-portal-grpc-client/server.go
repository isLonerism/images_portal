package main

import (
	"context"
	"log"
	"strconv"

	"github.com/bsuro10/images-portal/images-portal-grpc-server/api/docker"
	"google.golang.org/grpc"
)

func main() {
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	panic(err)
	// }

	// log.Print("Server started on port 8080")

	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := docker.NewDockerClient(conn)

	response, err := client.Load(context.Background(), &docker.S3Object{
		S3Key: "test-images.tar",
	})
	if err != nil {
		panic(err)
	}

	log.Println(response)

	var list []*docker.TagImage
	var counter int
	counter = 0
	for _, element := range response.GetImages() {
		counter = counter + 1
		list = append(list, &docker.TagImage{
			OldImage: element,
			NewImage: &docker.Image{
				Name: "bsuro10/tagggggeedddrantest" + strconv.Itoa(counter),
			},
		})
	}

	message, err := client.TagAndPush(context.Background(), &docker.TagAndPushObject{
		TagImages: &docker.TagImagesList{
			Images: list,
		},
		AuthConfig: &docker.AuthConfig{
			Username: "bsuro10",
			Password: "maginsuro",
		},
	})
	if err != nil {
		panic(err)
	}

	log.Println(message)
}
