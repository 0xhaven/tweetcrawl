package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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
	mux.HandleFunc("/info", rt.countHandler)
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
	var infoResp struct {
		Count int
		Duration string
		RatePerSecond float64
		PercentHashtag float64
		PercentURL float64
		PercentImageURL float64
		PercentEmoji float64
	}

	count, err := rt.store.Count()
	handleErr(w, err)
	duration, err := rt.store.Duration()
	handleErr(w, err)
	hashtagCount, err := rt.store.NumWithHashtags()
	handleErr(w, err)
	urlCount, err := rt.store.NumWithURLs()
	handleErr(w, err)
	imageCount, err := rt.store.NumWithPhotoURLs()
	handleErr(w, err)
	emojiCount, err := rt.store.NumWithEmoji()
	handleErr(w, err)

	infoResp.Count = count
	infoResp.Duration = duration.String()
	infoResp.RatePerSecond = float64(count) / float64(duration/time.Second)
	infoResp.PercentHashtag = float64(hashtagCount)/float64(count)
	infoResp.PercentURL = float64(urlCount)/float64(count)
	infoResp.PercentImageURL = float64(imageCount)/float64(count)
	infoResp.PercentEmoji = float64(emojiCount)/float64(count)
	handleErr(w, json.NewEncoder(w).Encode(infoResp))
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
