package main

import (
	"encoding/json"
	"fmt"
	"github.com/harry1453/go-common-file-dialog/cfd"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type Page struct {
	Title string
	Body  string
	Files []file
}

type file struct {
	Path string
	Name string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{
		Title: "File Rename Helper ;D",
	}

	t, err := template.ParseFiles("templates/index.html")
	ifErrorToPage(w, err)

	err = t.Execute(w, p)
	ifErrorToPage(w, err)
}

func main() {
	mux := http.NewServeMux()
	// libs
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	mux.Handle("/scss/", http.StripPrefix("/scss/", http.FileServer(http.Dir("./scss"))))
	mux.Handle("/vendor/", http.StripPrefix("/vendor/", http.FileServer(http.Dir("./vendor"))))
	// site
	mux.HandleFunc("/choose-files", chooseFilesHandler)
	mux.HandleFunc("/", indexHandler)
	_ = http.ListenAndServe(":8080", mux)
}

func ifErrorToPage(w io.Writer, err error) {
	if err != nil {
		t, e := template.ParseFiles("templates/Error.html")
		if e != nil {
			fmt.Println(e)
		}

		e = t.Execute(w, err)
		if e != nil {
			fmt.Println(e)
		}
	}
}

func chooseFilesHandler(w http.ResponseWriter, r *http.Request){
	filesJSON, err := json.Marshal(getFilesInDirectory())
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = w.Write(filesJSON)
	if err != nil {
		return
	}
}

func getFilesInDirectory()(files []file) {
	openMultiDialog, err := cfd.NewOpenMultipleFilesDialog(cfd.DialogConfig{
		Title:         "Open Multiple Files",
		Role:          "OpenFilesExample",
		FileFilters: []cfd.FileFilter{
			{
				DisplayName: "Text Files (*.txt)",
				Pattern:     "*.txt",
			},
			{
				DisplayName: "Image Files (*.jpg, *.png)",
				Pattern:     "*.jpg;*.png",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
		SelectedFileFilterIndex: 2,
		FileName:                "file.txt",
		DefaultExtension:        "txt",
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := openMultiDialog.Show(); err != nil {
		log.Fatal(err)
	}
	results, err := openMultiDialog.GetResults()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Chosen file(s): %s\n", results)

	for _, result := range results {
		files = append(files, file{
			Path: result,
			Name: getFileNameFromPath(result),
		})
	}

	return files
}

func getFileNameFromPath(path string)string{
	strArr := strings.Split(path,"\\")
	return strArr[len(strArr)-1]
}
}