package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", pageHandler)

	http.ListenAndServe(":8080", nil)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		panic(err)

	}

	var uploadedFiles []string

	for _, fh := range r.MultipartForm.File["file[]"] {

		file, _ := fh.Open()
		defer file.Close()

		out, err := os.Create("./uploadedFiles/" + fh.Filename)

		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing")
			return
		}

		defer out.Close()

		//write the content from POST to the file
		_, err = io.Copy(out, file)

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		uploadedFiles = append(uploadedFiles, fh.Filename)
	}

	if err := json.NewEncoder(w).Encode(uploadedFiles); err != nil {
		fmt.Fprintln(w, err)
		return
	}

}
