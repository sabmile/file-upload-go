package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const MAX_UPLOAD_FILE = 1024 * 1024 // 1 MB

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	http.ServeFile(w, r, "index.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_FILE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_FILE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose a file that's less than 1MB in size", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")

	defer file.Close()

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload succesful")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadHandler)

	if err := http.ListenAndServe(":4500", mux); err != nil {
		log.Fatal(err)
	}
}
