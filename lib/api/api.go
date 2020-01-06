package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jacobhaven/tweetcrawl/lib/store"
)

type router struct {
	http.Handler
	store store.Store
	mux *http.ServeMux
}

func NewRouter(store store.Store) http.Handler {
	mux := http.NewServeMux()
	rt := &router{
		Handler: mux,
		store: store,
	}
	mux.HandleFunc("/count", rt.countHandler)
	mux.HandleFunc("/hashtags", rt.topHashtags)
	mux.HandleFunc("/domains", rt.topDomains)
	mux.HandleFunc("/emoji", rt.topEmoji)
	return rt
}

func handleErr(w http.ResponseWriter, err error) {
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (rt *router) countHandler(w http.ResponseWriter, r *http.Request) {
	count, err := rt.store.Count()
	handleErr(w, err)
	handleErr(w, json.NewEncoder(w).Encode(struct{Count int}{count}))
}

func parseCount(val string) int {
	count, err := strconv.Atoi(val)
	if err != nil || count <= 0 || count > 100 {
		count = 5
	}
	return count
}

func (rt *router) topHashtags(w http.ResponseWriter, r *http.Request) {
	items, err := rt.store.TopHashtags(parseCount(r.URL.Query().Get("count")))
	handleErr(w, err)
	handleErr(w, json.NewEncoder(w).Encode(items))
}

func (rt *router) topDomains(w http.ResponseWriter, r *http.Request) {
	items, err := rt.store.TopDomains(parseCount(r.URL.Query().Get("count")))
	handleErr(w, err)
	handleErr(w, json.NewEncoder(w).Encode(items))
}

func (rt *router) topEmoji(w http.ResponseWriter, r *http.Request) {
	items, err := rt.store.TopEmoji(parseCount(r.URL.Query().Get("count")))
	handleErr(w, err)
	handleErr(w, json.NewEncoder(w).Encode(items))
}
