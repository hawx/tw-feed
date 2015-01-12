package main

import (
	"github.com/hawx/serve"
	"github.com/hawx/tw-feed/store"
	"github.com/hawx/tw-stream"
	"github.com/gorilla/feeds"

	"log"
	"flag"
	"fmt"
	"net/http"
	"time"
)

var (
	consumerKey    = flag.String("consumer-key", "", "")
	consumerSecret = flag.String("consumer-secret", "", "")
	accessToken    = flag.String("access-token", "", "")
	accessSecret   = flag.String("access-secret", "", "")

	port   = flag.String("port", "8080", "")
	socket = flag.String("socket", "", "")
	help   = flag.Bool("help", false, "")
)

const HELP = `Usage: tw-feed [options]

  Serves a rss feed of your twitter timeline.

    --consumer-key <value>
    --consumer-secret <value>
    --access-token <value>
    --access-secret <value>

    --port <port>       # Port to run on (default: '8080')
    --socket <path>     # Serve using a unix socket instead
    --help              # Display this help message
`

func main() {
	flag.Parse()

	if *help {
		fmt.Println(HELP)
		return
	}

	store := store.New(20)
	auth := stream.Auth(*consumerKey, *consumerSecret, *accessToken, *accessSecret)

	go func() {
		for tweet := range stream.Self(auth) {
			store.Add(tweet)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tweets := store.Latest()
		if len(tweets) == 0 {
			w.WriteHeader(204)
			return
		}

		feed := &feeds.Feed{
			Title:   tweets[0].User.Name,
			Link:    &feeds.Link{Href: *tweets[0].User.Url},
		  Created: time.Now(),
		}

		for _, tweet := range tweets {
			feed.Items = append(feed.Items, &feeds.Item{
			  Link:        &feeds.Link{Href: tweet.Link()},
				Description: tweet.Text,
				Created:     tweet.CreatedAt.Time,
			})
		}

		w.Header().Add("Content-Type", "application/rss+xml")

		err := feed.WriteRss(w)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
	})

	serve.Serve(*port, *socket, http.DefaultServeMux)
}
