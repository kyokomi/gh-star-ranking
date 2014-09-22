package ghstar

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"appengine"
	"appengine/datastore"
)

const dataStoreRanking = "Ranking"

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

	// TODO: 前日のランキングを取得
//	rankings := make([]Ranking, 0, 30)
//	if _, err := q.GetAll(c, &rankings); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}

	if err := rankingTemplate.Execute(w, res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	vars := mux.Vars(r)
	lang := vars["language"]

	// TODO: 今日のデータを取得するようにする
	q := datastore.NewQuery(dataStoreRanking).Ancestor(rankingKey(c)).Filter("Lang = ", lang).Order("-Date").Limit(30)
	count, err := q.Count(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 今日の登録
	if count == 0 {
		rankings, err := readGitHubStarRanking(c, lang)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, ranking := range rankings {
			key := datastore.NewIncompleteKey(c, dataStoreRanking, rankingKey(c))
			_, err := datastore.Put(c, key, &ranking)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		fmt.Fprint(w, "登録した ", len(rankings))
	} else {
		fmt.Fprint(w, "登録済み ")
	}
}

func rankingKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, dataStoreRanking, "defualt_ranking", 0, nil)
}
