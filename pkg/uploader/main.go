package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

// CREDIT: https://github.com/TannerGabriel/learning-go/tree/master/beginner-programs/FileUpload

// // Compile templates on start of the application
var templates = template.Must(template.ParseFiles("public/upload.html"))

// // Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}

	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create(handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func formDataFileUploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}

func byteStreamFileUploadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.Header.Get("X-Filename")
	folder := r.Header.Get("X-Folder")
	filepath := folder + "/" + filename

	log.Println("destination: " + filepath)

	file, err := os.Create(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	n, err := io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := []byte(fmt.Sprintf("%d bytes are recieved.\n", n))
	log.Println(string(response))

	w.Write(response)
}

func main() {
	http.HandleFunc("/file-upload", formDataFileUploadHandler)
	http.HandleFunc("/stream-upload", byteStreamFileUploadHandler)

	//Listen on port 8080
	log.Println("starting uploader...")
	http.ListenAndServe(":8080", nil)
}
