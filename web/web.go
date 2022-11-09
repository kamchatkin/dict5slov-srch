package web

import (
	_ "embed"
	"encoding/json"
	"kamchatkin.ru/wordle-hack/search"
	"net/http"
)

//go:embed index.html
var IndexPage string

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/kamchatkin/wordle-hack", http.StatusTemporaryRedirect)

	//_, err := fmt.Fprintf(w, IndexPage)
	//if err != nil {
	//	log.Panic(err)
	//}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	params := search.QueryConstructor(
		q.Get("letter0"),
		q.Get("ne0"),
		q.Get("letter1"),
		q.Get("ne1"),
		q.Get("letter2"),
		q.Get("ne2"),
		q.Get("letter3"),
		q.Get("ne3"),
		q.Get("letter4"),
		q.Get("ne44"),
		q.Get("lettersI"),
		q.Get("lettersE"))

	words := search.WebSearch(params)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(words)
}

func Web() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	_ = http.ListenAndServe(":8080", nil)
}
