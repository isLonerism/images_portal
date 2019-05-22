package main

import (
	"log"
	"net/http"

	oc "github.com/bsuro10/images-portal/images-portal-grpc-client/interfaces/openshift_client"
)

func main() {
	http.HandleFunc("/upload", uploadFileHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

	log.Print("Server started on port 8080")
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	oc.DeployPod()
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
