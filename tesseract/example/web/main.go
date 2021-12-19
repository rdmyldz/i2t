package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/rdmyldz/i2t/tesseract"
)

// Compile templates on start of the application
var templates = template.Must(template.ParseFiles("index.html"))

// Display the named template
func (app *app) display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func (app *app) uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	f, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("error reading r.Body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("fsize in handler: %v\n", len(f))
	texts, err := app.handle.ProcessImageMem(f)
	if err != nil {
		log.Println("error ProcessImageMem: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*
		// app.handle.SetImage2(handler.Filename)
		app.handle.SetImage2FromMem(f)
		app.handle.Recognize()
		text, _ := app.handle.GetUTF8Text()

		// app.handle.End()
		// app.handle.Delete()
	*/

	fmt.Fprint(w, texts)
}

func (app *app) uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.display(w, "index", nil)
	case "POST":
		app.uploadFile(w, r)
	}
}

type app struct {
	handle *tesseract.TessBaseAPI
}

func main() {

	handle, err := tesseract.TessBaseAPICreate("tur+eng")
	if err != nil {
		log.Fatalln(err)
	}

	app := app{
		handle: handle,
	}
	// Upload route
	http.HandleFunc("/upload", app.uploadHandler)

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
