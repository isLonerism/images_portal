package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

const maxUploadSize = 1 * 1024 * 1024 // 1 GB
const uploadPath = "./tmp"

func main() {
	http.HandleFunc("/upload", uploadFileHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

	log.Print("Server started on port 8080")
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: why do we need MaxBytesReader and ParseMultipartForm both
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		log.Println(err)
		return
	}

	fileType := r.PostFormValue("type")
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	defer file.Close() // closing the file after the function is done
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
