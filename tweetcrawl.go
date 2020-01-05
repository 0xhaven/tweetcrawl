package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jacobhaven/tweetcrawl/lib/twitter"
)

func main() {
	sampler,  err := twitter.NewSampler(
		os.Getenv("CUSTOMER_KEY"),
		os.Getenv("CUSTOMER_SECRET"),
		)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := sampler.Open()
	if err !=nil {
		log.Fatal(err)
	}

	go func() {
		for tweet := range stream {
			fmt.Println(tweet)
		}
	}()
	time.Sleep(time.Second)
	if err := sampler.Close(); err != nil {
		log.Fatal(err)
	}
}