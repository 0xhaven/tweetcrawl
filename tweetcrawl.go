package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jacobhaven/tweetcrawl/lib/api"
	"github.com/jacobhaven/tweetcrawl/lib/store"
	"github.com/jacobhaven/tweetcrawl/lib/twitter"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

)

func main() {
	db, err := gorm.Open("sqlite3","tweet.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store, err := store.NewSQL(db)
	if err != nil {
		log.Fatal(err)
	}

	sampler,  err := twitter.NewSampler(
		os.Getenv("CUSTOMER_KEY"),
		os.Getenv("CUSTOMER_SECRET"),
	)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := sampler.Open()
	if err != nil {
		log.Fatal(err)
	}

	const numWorkers = 100

	for i := 0; i < numWorkers; i++ {
		go func() {
			for tweet := range stream {
				store.Save(tweet)
			}
		}()
	}

	const addr = ":8080"
	log.Printf("Listing on %s\n", addr)
	if err := http.ListenAndServe(addr, api.NewRouter(store)); err != nil {
		log.Fatal(err)
	}
}