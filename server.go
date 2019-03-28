package main

import (
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
	maxUploadSize = 1 * 1024 * 1024 // 1 GB
	s3_bucket     = "uploaded-images"
)

var (
	s3_config = &aws.Config{
		Credentials:      credentials.NewStaticCredentials("FZJ94C79V4MTQUVIESCC", "c6aF+ZELuUnd8CgL0wnFnSFoSXE2pjwNeTJyQIg0", ""),
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
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		log.Println(err)
		return
	}

	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		renderError(w, "ERROR_READING_FILE_HEADER", http.StatusBadRequest)
		log.Println(err)
		return
	}

	filetype, err := filetype.Match(fileHeader)
	if err != nil {
		renderError(w, "SOMETHING_WENT_WRONG_WITH_YOUR_FILE_TYPE", http.StatusBadRequest)
		log.Println(err)
		return
	}
	if filetype.Extension != "tar" {
		renderError(w, "NOT_SUPPORTED_FILE_TYPE", http.StatusBadRequest)
		fmt.Println("File type is not supported: ", filetype.Extension)
		return
	}

	fmt.Println(filetype.Extension)

	defer file.Close()

	result, err := uploadToS3(&file, handler.Filename)
	if err != nil {
		renderError(w, "UPLOAD_FAILED", http.StatusBadRequest)
		log.Println(err)
		return
	}

	fmt.Printf("File uploaded to: %s\n", aws.StringValue(&result.Location))
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

func validateFile() {

}

func pushToOpenshift() {

}
