package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/harry1453/go-common-file-dialog/cfd"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type Page struct {
	Title string
	Body  string
	Files []file
}

type file struct {
	UUID          string
	Path          string
	TruncatedPath string
	Name          string
	Extension     string
	NewName       string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{
		Title: "file Rename Helper ;D",
	}

	t, err := template.ParseFiles("templates/index.html")
	ifErrorToPage(w, err)

	err = t.Execute(w, p)
	ifErrorToPage(w, err)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/choose-files", chooseFilesHandler).Methods("GET")
	r.HandleFunc("/rename-selected-files", renameSelectedFilesHandler).Methods("POST")
	
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	http.Handle("/", r)
	_ = http.ListenAndServe(":8080", nil)
}

func renameSelectedFilesHandler(w http.ResponseWriter, r *http.Request) {
	var selectedFiles []file

	if err := json.NewDecoder(r.Body).Decode(&selectedFiles);err != nil {
		http.Error(w, err.Error(), 400)
		GetLogger(r).WithError(err)

		return
	}

	fmt.Println(selectedFiles)
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

func chooseFilesHandler(w http.ResponseWriter, r *http.Request) {
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

func getFilesInDirectory() (files []file) {
	openMultiDialog, err := cfd.NewOpenMultipleFilesDialog(cfd.DialogConfig{
		Title: "Open Multiple Files",
		Role:  "OpenFilesExample",
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
		fmt.Println(err)

		return nil
	}

	log.Printf("Chosen file(s): %s\n", results)

	for _, result := range results {
		files = append(files, file{
			Path:          result,
			TruncatedPath: truncatePath(result),
			Name:          getFileNameFromPath(result),
			Extension:     getExtention(result),
		})
	}

	return files
}

func getFileNameFromPath(path string) string {
	strArr := strings.Split(path, "\\")

	return strArr[len(strArr)-1]
}

func truncatePath(path string) string {
	return strings.TrimSuffix(path, getFileNameFromPath(path))
}

func getExtention(path string) string {
	strArr := strings.Split(path, ".")

	return strArr[len(strArr)-1]
}

func GetLogger(r *http.Request) *log.Entry {
	return r.Context().Value("logger").(*log.Entry)
}
