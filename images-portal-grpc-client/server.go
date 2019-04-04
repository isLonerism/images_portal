package main

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/h2non/filetype"
)

// TODO: Pass to enviornment variables
const (
	maxUploadSize = 10 * 1024 * 1024 // 10 GB
	s3_bucket     = "uploaded-images"
)

var (
	s3_config = &aws.Config{
		Credentials:      credentials.NewStaticCredentials("EUZWJJWSWHS3TWQ0MQKN", "OcJJ1NIEGozBx0yvmg3oXZrxE+3k2T2e2B3N3cHZ", ""),
		Endpoint:         aws.String("http://localhost:9000"),
		Region:           aws.String("default"),
		S3ForcePathStyle: aws.Bool(true),
	}
)

func main() {
	http.HandleFunc("/upload", uploadFileHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

	log.Print("Server started on port 8080")
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w, "FILE_TOO_BIG/INVALID_FORM_DATA", http.StatusBadRequest)
			log.Println(err)
			return
		}

		file, handler, err := r.FormFile("uploadFile")
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		if err := validateFile(&file); err != nil {
			renderError(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		}

		defer file.Close()

		result, err := uploadToS3(&file, handler.Filename)
		if err != nil {
			renderError(w, "UPLOAD_FAILED", http.StatusBadRequest)
			log.Println(err)
			return
		}

		fmt.Printf("File uploaded to: %s\n", aws.StringValue(&result.Location))
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func uploadToS3(file *multipart.File, filename string) (*s3manager.UploadOutput, error) {
	sess := session.Must(session.NewSession(s3_config))
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3_bucket),
		Key:    aws.String(filename),
		Body:   *file,
	})

	return result, err
}

func validateFile(file *multipart.File) error {
	fileHeader := make([]byte, 512)
	if _, err := (*file).Read(fileHeader); err != nil {
		return err
	}

	filetype, err := filetype.Match(fileHeader)
	if err != nil {
		return err
	}
	if filetype.Extension != "tar" {
		return errors.New("File type is not supported: " + filetype.Extension)
	}

	// TODO: Validate filesize

	return nil
}

func pushToOpenshift() {

}
