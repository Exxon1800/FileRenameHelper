package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  string
	Files []file
}

type file struct {
	Name string
	Path string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{
		Title: "File Rename Helper ;D",
		Body:  "hi",
		Files: []file{{"file1", ""}, {"file2", ""}, {"file3", ""}},
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
	mux.HandleFunc("/submit-files", FileHandler)
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

func FileHandler(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r.PostForm)
}

func getFilesInDirectory(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
