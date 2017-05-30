package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", defaultHandler)
	r.HandleFunc("/courses", coursesHandler).Methods("GET")
	r.HandleFunc("/course/{name}/menu", courseMenuHandler).Methods("GET")
	r.HandleFunc("/course/{name}/{chapter}", courseContentHandler).Methods("GET")
	http.Handle("/", r)
}

func coursesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"courses":[{"name":"Foo","slug":"foo"}]}`))
}

func courseMenuHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseName := vars["name"]

	ctx := appengine.NewContext(r)
	content, err := readFileFromBucket(ctx, "courses/"+courseName+"/_menu.json")
	if err != nil {
		log.Errorf(ctx, "could not read _menu.json of %s: %v", courseName, err)
		w.WriteHeader(http.StatusTeapot)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

func courseContentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseName := vars["name"]
	chapterName := vars["chapter"]

	ctx := appengine.NewContext(r)
	content, err := readFileFromBucket(ctx, "courses/"+courseName+"/"+chapterName+".html")
	if err != nil {
		log.Errorf(ctx, "could not read %s chapter %s: %v", courseName, chapterName, err)
		w.WriteHeader(http.StatusTeapot)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(content)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
