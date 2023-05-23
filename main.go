package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)           // limit your max memory size
	file, handler, err := r.FormFile("file") // this parses the multipart form
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Define the directory
	dir := "./uploads"

	// Check if the directory exists, if not, create it
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Create the file
	dst, err := os.Create("./uploads/" + handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully!"))
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	// Assume you want to get the file name from URL like "/download/filename"
	fileName := r.URL.Path[len("/download/"):]

	http.ServeFile(w, r, "./uploads/"+fileName)
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	dir := "./uploads"

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]")) // No files, return empty array
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(files) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]")) // No files, return empty array
		return
	}

	var fileURLs []string
	for _, file := range files {
		if !file.IsDir() {
			// Build full URL
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			} else if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
				scheme = forwardedProto
			}
			fileURL := fmt.Sprintf("%s://%s/download/%s", scheme, r.Host, file.Name())
			fileURLs = append(fileURLs, fileURL)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(fileURLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`404: file is not found

Try to upload a file using the following command:

upload: 

  curl -v -F "file=@/path/to/file.ext" http://localhost:3000/upload
  curl -v -F "file=@/path/to/some.zip;filename=fancy.zip" http://localhost:3000/upload

download:

  curl -O http://localhost:3000/download/file.ext
  curl -O http://localhost:3000/download/fancy.zip

list files:

  curl http://localhost:3000/files

Compress files:

  tar -cf compressed.tar /path/to/directory
  tar -cf - /path/to/directory | pigz > compressed.tar.gz
  tar -czf compressed.tar.gz /path/to/directory

Decompress files:

  tar -xf compressed.tar
  tar -xzf compressed.tar.gz
`))
}

func main() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download/", downloadFile)
	http.HandleFunc("/files", listFiles)
	http.HandleFunc("/", notFoundHandler)

	fmt.Println("Server starting and listening on :3000...")
	http.ListenAndServe(":3000", nil)
}
