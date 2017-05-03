package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
	"time"
)

// constants
const (
	twitterTimeStringFmt = "15:04 PM - 2 Jan 2006"
	localTimeZone = "CST" // timezone where the compute is happening
)

type TweetHeader struct {
	tweetUrl   string
	time       int64 // unix time
	stringTime string
}
type TweetFooter struct {
	favorites int
	retweets  int
}
type Tweet struct {
	tweetText   string
	hashTags    []string
	tweetHeader TweetHeader
	tweetFooter TweetFooter
}

func main() {
	// this program demonstrates scraping
	// TODO:
	// 	- check if private account of not
	// 	- continuous loading (lookup api req to load more)
	// 	- set timezone
	// 	- think about data stores?
	// 	- finish tweetfooter retriever, will probably need to follow actual tweet link to get
	//	  all the likes, retweets, and replies available for public (big ticket here)

	resp, err := http.Get("https://twitter.com/realDonaldTrump")
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	tweets := []Tweet{}
	var currTweet Tweet
	currEleTweet := false

	// tree traversal
	var htmlTweetParser func(*html.Node)
	htmlTweetParser = func(n *html.Node) {
		// check if node type is div, if yes investigate attributes of the node
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" {
					if a.Val == "stream-item-header" {
						currTweet.tweetHeader = retrieveTweetHeader(n)
						currEleTweet = true
					} else if a.Val == "js-tweet-text-container" {
						currTweet.tweetText = retrieveTweet(n)
					} else if a.Val == "stream-item-footer" {
						currTweet.tweetFooter = retrieveTweetFooter(n)
					}
				}
			}
			// if current node contains attributes that signal a tweet, add to
			// tweet slice and reset currEleTweet flag
			if currEleTweet{
				currEleTweet = false
				tweets = append(tweets, currTweet)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlTweetParser(c)
		}
	}
	htmlTweetParser(root)
	fmt.Println(tweets)
}

// retrieve 
func retrieveTweet(n *html.Node) string {
	tweet := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tweet += retrieveTweet(c)
	}
	if n.Type == html.TextNode {
		return string(n.Data)
	}
	return tweet
}

func retrieveTweetHeader(n *html.Node) TweetHeader {
	var header TweetHeader

	var processTweetHeaderHelper func(n *html.Node)
	processTweetHeaderHelper = func(n *html.Node) {
		for e := n.FirstChild; e != nil; e = e.NextSibling {
			for _, a1 := range e.Attr {
				if a1.Key == "title" {
					header.stringTime = a1.Val
					header.time = stringTimeToUnixTime(header.stringTime)
				}
				if a1.Key == "href" {
					header.tweetUrl = a1.Val
				}
			}

		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		for _, a := range c.Attr {
			if a.Val == "time" {
				processTweetHeaderHelper(c)
			}
		}
	}
	return header
}

func retrieveTweetFooter(n *html.Node) TweetFooter {
	var footer TweetFooter
	var processTweetFooterHelper func(n *html.Node)
	processTweetFooterHelper = func(n *html.Node) {
		for e := n.FirstChild; e != nil; e = e.NextSibling {
			for _, a1 := range e.Attr {
				if a1.Key == "title" {
					fmt.Println("time: ", a1.Val)
				}
				if a1.Key == "href" {
					fmt.Println("link: ", a1.Val)
				}
			}

		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		for _, a := range c.Attr {
			if a.Val == "time" {
				processTweetFooterHelper(c)
			}
		}
	}
	return footer
}

// perhaps split time operations into its own util file
// https://golang.org/src/time/format.go <--- for reference
func stringTimeToUnixTime(stringTime string) int64 {
	tweetTime, err := time.Parse(twitterTimeStringFmt, stringTime)
	if err != nil {
		return -1
	}
	return tweetTime.Unix()
}


