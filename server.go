package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	oc "github.com/bsuro10/images-portal/interfaces/openshift_client"
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

	// TODO: Add fileType check (only .tar and .tar.gz)

	// Upload the file to s3 object storage
	sess := session.Must(session.NewSession(s3_config))
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3_bucket),
		Key:    aws.String(handler.Filename),
		Body:   file,
	})

	if err != nil {
		renderError(w, "UPLOAD_FAILED", http.StatusBadRequest)
		log.Println(err)
		return
	}

	defer file.Close()

	oc.DeployPod()

	fmt.Printf("File uploaded to: %s\n", aws.StringValue(&result.Location))
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
