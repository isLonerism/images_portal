package docker

import (
	"bytes"
	context "context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	s3_bucket = "uploaded-images"
)

var (
	s3_config = &aws.Config{
		Credentials:      credentials.NewStaticCredentials("7604HNT325T2DMRRP7T1", "ZhpOcDA5kTr8gBA0+7pq7DeCS7AdeqAA4krotRO3", ""),
		Endpoint:         aws.String("http://localhost:9000"),
		Region:           aws.String("default"),
		S3ForcePathStyle: aws.Bool(true),
	}
)

type DockerServiceServer struct {
	client *client.Client
}

func RegisterDockerService(client *client.Client) *DockerServiceServer {
	return &DockerServiceServer{client: client}
}

func (s *DockerServiceServer) Load(ctx context.Context, s3_object *S3Object) (*ImagesList, error) {
	buff, err := downloadS3file((*s3_object).GetS3Key(), s3_bucket, s3_config)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error occurred while downloading the image")
	}

	response, err := s.client.ImageLoad(ctx, bytes.NewReader((*buff).Bytes()), true)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error occurred while loading the images")
	}
	defer response.Body.Close()

	images, err := getImagesList(response)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error occurred while getting the images")
	}

	var imagesList []*Image
	for _, element := range *images {
		imagesList = append(imagesList, &Image{
			Name: element,
		})
	}

	return &ImagesList{
		Images: imagesList,
	}, nil
}

func downloadS3file(key string, bucket string, config *aws.Config) (*aws.WriteAtBuffer, error) {
	buff := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(session.New(config))
	_, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return buff, err
}

func getImagesList(response types.ImageLoadResponse) (*[]string, error) {
	images, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	reg := regexp.MustCompile(`Loaded image: ([^\s\\]+)`)
	list := reg.FindAllString(string(images), -1)

	for index, element := range list {
		list[index] = strings.Fields(element)[2]
	}

	return &list, nil
}

func (s *DockerServiceServer) TagAndPush(ctx context.Context, tagAndPushObject *TagAndPushObject) (*Message, error) {
	for _, element := range (*tagAndPushObject).GetTagImages().GetImages() {

		old_image := element.GetOldImage().GetName()
		new_image := element.GetNewImage().GetName()
		if err := s.client.ImageTag(ctx, old_image, new_image); err != nil {
			log.Println(err)
			return nil, errors.New("Error occurred while tagging the images")
		}

		registryAuth, err := encodedAuthConfig((*tagAndPushObject).GetAuthConfig().GetUsername(),
			(*tagAndPushObject).GetAuthConfig().GetPassword())
		if err != nil {
			log.Println(err)
			return nil, errors.New("Error occurred while encoding the auth config")
		}

		_, err = s.client.ImagePush(ctx, new_image, types.ImagePushOptions{
			RegistryAuth: registryAuth,
		})
		if err != nil {
			log.Println(err)
			return nil, errors.New("Error occurred while pushing the image: " + new_image)
		}
	}

	return &Message{
		Message: "Successfull pushed!",
	}, nil
}

func encodedAuthConfig(username string, password string) (string, error) {
	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encodedJSON), nil
}
