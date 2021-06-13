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
	"os"
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
	ifErrorToPage(w, r, err)

	err = t.Execute(w, p)
	ifErrorToPage(w, r, err)
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

func renameSelectedFilesHandler(_ http.ResponseWriter, r *http.Request) {
	var selectedFiles []file

	err := json.NewDecoder(r.Body).Decode(&selectedFiles)
	if logIfError(r, err){
		return
	}

	err = renameFiles(selectedFiles)
	if logIfError(r, err){
		return
	}

	log.Printf("Selected file(s): %s\n", selectedFiles)
}

func chooseFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := getFilesInDirectory()
	if logIfError(r, err) {
		return
	}

	filesJSON, err := json.Marshal(files)
	if logIfError(r, err) {
		return
	}

	_, err = w.Write(filesJSON)
	if err != nil {
		return
	}
}

func getFilesInDirectory() (files []file, err error) {
	var numberOfFileTypeFilters uint = 2

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
		SelectedFileFilterIndex: numberOfFileTypeFilters,
		FileName:                "file.txt",
		DefaultExtension:        "txt",
	})
	if err != nil {
		return nil,  fmt.Errorf("could not create NewOpenMultipleFilesDialog %w", err)
	}

	if err = openMultiDialog.Show(); err != nil {
		return nil,  fmt.Errorf("could not show openMultiDialog %w", err)
	}

	results, err := openMultiDialog.GetResults()
	if err != nil {
		return nil,  fmt.Errorf("could not get results from openMultiDialog %w", err)
	}

	log.Printf("Chosen file(s): %s\n", results)

	for _, result := range results {
		files = append(files, file{
			Path:          result,
			TruncatedPath: truncatePath(result),
			Name:          getFileNameFromPath(result),
			Extension:     getExtension(result),
		})
	}

	return files, nil
}

func renameFiles(files []file) error {
	for _, file := range files {
		newFilePath := file.TruncatedPath + file.NewName
		err := os.Rename(file.Path, newFilePath)

		if err != nil {
			return fmt.Errorf("could not rename file %w", err)
		}
	}

	return nil
}

func getFileNameFromPath(path string) string {
	strArr := strings.Split(path, "\\")

	return strArr[len(strArr)-1]
}

func truncatePath(path string) string {
	return strings.TrimSuffix(path, getFileNameFromPath(path))
}

func getExtension(path string) string {
	strArr := strings.Split(path, ".")

	return strArr[len(strArr)-1]
}

func GetLogger(r *http.Request) *log.Entry {
	return r.Context().Value("logger").(*log.Entry)
}

func logError(r *http.Request, err error) *log.Entry {
	return GetLogger(r).WithError(err)
}

func logIfError(r *http.Request, err error) bool {
	if err != nil {
		logError(r, err)

		return true
	}

	return false
}

func ifErrorToPage(w io.Writer, r *http.Request, err error) {
	if err != nil {
		t, err := template.ParseFiles("templates/Error.html")
		if logIfError(r, err) {
			return
		}

		err = t.Execute(w, err)
		if logIfError(r, err) {
			return
		}
	}
}
