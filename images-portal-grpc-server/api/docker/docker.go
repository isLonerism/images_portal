package docker

import (
	"bytes"
	context "context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
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

type DockerServiceServer struct {
	client *client.Client
}

func RegisterDockerService(client *client.Client) *DockerServiceServer {
	return &DockerServiceServer{client: client}
}

func (s *DockerServiceServer) Load(ctx context.Context, s3_object *S3Object) (*ImagesList, error) {
	buff, err := downloadS3file((*s3_object).GetS3Key(), (*s3_object).GetS3Bucket(), &aws.Config{
		Credentials:      credentials.NewStaticCredentials((*s3_object).GetS3Accesskey(), (*s3_object).GetS3Secretkey(), ""),
		Endpoint:         aws.String((*s3_object).GetS3Endpoint()),
		Region:           aws.String((*s3_object).GetS3Region()),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error occurred while downloading the image: " + err.Error())
	}

	response, err := s.client.ImageLoad(ctx, bytes.NewReader((*buff).Bytes()), true)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error occurred while loading the images: " + err.Error())
	}
	defer response.Body.Close()

	imagesList, err := convertImageLoadResponseToImagesList(response)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error occurred while converting the images: " + err.Error())
	}

	return &ImagesList{
		Images: imagesList,
	}, nil
}

func convertImageLoadResponseToImagesList(response types.ImageLoadResponse) ([]*Image, error) {
	images, err := getImagesList(response)
	if err != nil {
		return nil, err
	}

	var imagesList []*Image
	for _, element := range *images {
		imagesList = append(imagesList, &Image{
			Name: element,
		})
	}

	return imagesList, nil
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
			return nil, errors.New("Error occurred while tagging the images: " + err.Error())
		}

		registryAuth, err := encodedAuthConfig((*tagAndPushObject).GetAuthConfig().GetUsername(),
			(*tagAndPushObject).GetAuthConfig().GetPassword())
		if err != nil {
			log.Println(err)
			return nil, errors.New("Error occurred while encoding the auth config: " + err.Error())
		}

		err = pushImage(s, ctx, registryAuth, new_image)
		if err != nil {
			log.Println(err)
			return nil, errors.New("Error occurd while pushing the images: " + err.Error())
		}
	}

	return &Message{
		Message: "Successfull pushed!",
	}, nil
}

func pushImage(s *DockerServiceServer, ctx context.Context, auth string, image string) error {

	out, err := s.client.ImagePush(ctx, image, types.ImagePushOptions{
		RegistryAuth: auth,
	})
	if err != nil {
		return err
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	return nil
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
