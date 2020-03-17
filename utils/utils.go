package utils

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const (
	PathToImagesDir = "img"
	PathFileServer  = "static"
)

func FileUpload(r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20) // Max memory 32 MiB
	file, handler, err := r.FormFile("picture")
	if err != nil {
		log.Printf("Error getting file: %s", err.Error())
		return "", err
	}
	defer file.Close() // Close the file when finished

	newFileName, err := generateRandomFilename(PathToImagesDir, handler.Filename)
	if err != nil {
		log.Printf("Error creating file: %s", err.Error())
		return "", err
	}
	f, _ := os.OpenFile(path.Join(PathToImagesDir, newFileName), os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	io.Copy(f, file)
	return path.Join(PathFileServer, newFileName), nil
}

func ResponseJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func CheckErr(err error) {
	if err != nil {
		log.Print(err)
	}
}

func generateRandomFilename(dir, filename string) (name string, err error) {
	for i := 0; i < 10000; i++ {
		/*
			Careful below! From documentation on rand:
			Random numbers are generated by a Source. Top-level functions, such as Float64 and Int, use a default shared Source that produces a DETERMINISTIC sequence of values each time a program is run. Use the Seed function to initialize the default Source if different behavior is required for each run.
		*/
		name = strconv.Itoa(rand.Int()) + path.Ext(filename)
		namePath := filepath.Join(dir, name)
		_, err := os.OpenFile(namePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		defer os.Remove(namePath)
		if os.IsExist(err) {
			continue
		}
		break
	}
	return
}
