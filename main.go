package ghstar

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"appengine"
)

func init() {
	// router
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
	r.HandleFunc("/ranking/{language}", rankingHandler).Methods("GET")
	r.HandleFunc("/snapshot/{language}", snapshotHandler).Methods("GET")
	// static
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))
	// handle
	http.Handle("/", r)
}

var rankingTemplate = template.Must(template.ParseFiles("templates/ranking.tmpl"))

func handler(w http.ResponseWriter, r *http.Request) {
	appengine.NewContext(r)

	if err := rankingTemplate.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rankingHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	vars := mux.Vars(r)
	lang := vars["language"]

	rankings, err := readGitHubStarRanking(c, lang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ResponseRanking{
		Language: lang,
		Rankings: rankings,
	}
	//	c.Infof(fmt.Sprint(res))

	if err := rankingTemplate.Execute(w, res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
