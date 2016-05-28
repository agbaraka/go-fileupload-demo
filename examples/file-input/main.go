package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	out, err := os.OpenFile("./uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing")
		return
	}

	defer out.Close()

	//write the content from POST to the file
	_, err = io.Copy(out, file)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintf(w, header.Filename)

}
