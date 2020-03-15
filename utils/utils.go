package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const pathToImagesDir = "img/"

func FileUpload(r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20) // Max memory 32 MiB
	file, handler, err := r.FormFile("picture")
	if err != nil {
		log.Printf("Error getting file: %s", err.Error())
		return "", err
	}
	defer file.Close() // Close the file when finished

	fullPath := pathToImagesDir + handler.Filename
	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
		return "", err
	}
	defer f.Close()

	io.Copy(f, file)
	return fullPath, nil
}

func ResponseJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
