package web

import (
	_ "embed"
	"encoding/json"
	"kamchatkin.ru/wordle-hack/search"
	"net/http"
)

// go:embed index.html
//var IndexPage string

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/kamchatkin/wordle-hack", http.StatusTemporaryRedirect)

	//_, err := fmt.Fprintf(w, IndexPage)
	//if err != nil {
	//	log.Panic(err)
	//}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	query, err := search.QueryConstructor(
		q.Get("letter0"),
		q.Get("letter1"),
		q.Get("letter2"),
		q.Get("letter3"),
		q.Get("letter4"),
		q.Get("ne0"),
		q.Get("ne1"),
		q.Get("ne2"),
		q.Get("ne3"),
		q.Get("ne4"),
		q.Get("lettersE"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	words := search.WebSearch(query)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":    true,
		"words": words,
	})
}

func Web() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	_ = http.ListenAndServe(":8080", nil)
}
