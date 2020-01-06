package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	time.Sleep(10*time.Second)
	if err := sampler.Close(); err != nil {
		log.Fatal(err)
	}

	count, err := store.Count()
	if err != nil {
		log.Fatal(err)
	}

	hashtags, err := store.TopHashtags(5)
	if err != nil {
		log.Fatal(err)
	}

	domains, err := store.TopDomains(7)
	if err != nil {
		log.Fatal(err)
	}

	emoji, err := store.TopEmoji(10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(count)
	for _, item := range hashtags {
		fmt.Println(item.Count, item.Name)
	}
	for _, item := range domains {
		fmt.Println(item.Count, item.Name)
	}
	for _, item := range emoji {
		fmt.Println(item.Count, item.Name)
	}
}