package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
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
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" {
					if a.Val == "stream-item-header" {
						currTweet.tweetHeader = processTweetHeader(n)
						currEleTweet = true
					} else if a.Val == "js-tweet-text-container" {
						currTweet.tweetText = processTweet(n)
					}
				}
			}
			if currEleTweet{
				currEleTweet = false
				tweets = append(tweets, currTweet)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(root)
	fmt.Println(tweets)
}

func processTweet(n *html.Node) string {
	tweet := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tweet += processTweet(c)
	}
	if n.Type == html.TextNode {
		return string(n.Data)
	}
	return tweet
}

func processTweetHeader(n *html.Node) TweetHeader {
	var header TweetHeader

	var processTweetHeaderHelper func(n *html.Node)
	processTweetHeaderHelper = func(n *html.Node) {
		for e := n.FirstChild; e != nil; e = e.NextSibling {
			for _, a1 := range e.Attr {
				if a1.Key == "title" {
					fmt.Println("time: ", a1.Val)
					header.stringTime = a1.Val
					header.time = stringTimeToUnixTime(header.stringTime)
				}
				if a1.Key == "href" {
					fmt.Println("link: ", a1.Val)
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

// TODO this guy
func processTweetFooter(n *html.Node) TweetFooter {
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

func stringTimeToUnixTime(stringTime string) int64 {
	return 1
}


func insertTweet() {

}
