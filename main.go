package main

import (
	"fmt"
	"flag"
)

func main() {
	// TODO 1: change all printf to log statements w proper logID
	// TODO 2: dynamic page loading
	// TODO 3: page loading delay to avoid twitter's scraping detection
	// TODO 4: (more devOps-y type shit) vpn switching to fly under the radar
	// TODO 5: setup GRPC messaging
	// TODO 6: Postgres integration
	// TODO 7: schedule high profile account jobs
	// TODO 8: setup influxdb to track stats like followers, followings, follower ratio, etc

	twitterHandlePtr := flag.String("handle", "", "please supply twitter handle")
	numTweetsPtr := flag.Int("num_tweets", 20, "please supply number of tweets to process")
	flag.Parse()
	twitterController(*twitterHandlePtr, *numTweetsPtr)
}

func twitterController(twitterHandle string, numTweets int) {
	tweets := getTweets(twitterHandle)
	if len(tweets) == 0 {
		fmt.Println(fmt.Sprintf("No tweets for accout: %s", twitterHandle))
		return
	} else {
		fmt.Println(tweets)
	}
	for len(tweets) < numTweets {
		tweets = append(tweets, getTweets(twitterHandle)...)
		// TODO put delay in here so we don't trip twitter's fkn bullshit
	}

	fmt.Println(tweets)
}