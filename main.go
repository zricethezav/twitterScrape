package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
	//"time"
	// "strings"
)

type Tweet struct {
	tweetText string
	hashTags  []string
	tweetUrl  string
	time      float64 // unix time
}

func main() {
	// this program demonstrates scraping
	// TODO:
	// 	- check if private account of not
	//	- get actual status link
	// 	- continuous loading (lookup api req to load more)

	resp, err := http.Get("https://twitter.com/realDonaldTrump")
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	// tweitter container
	//var tweets []Tweet

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" {
					if a.Val == "stream-item-header" {
						processTweetHeader(n)
						break
					}
					if a.Val == "js-tweet-text-container" {
						processTweet(n)
						break
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(root)
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

func processTweetHeader(n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		for _, a := range c.Attr {
			if a.Val == "time" {
				for e := c.FirstChild; e != nil; e = e.NextSibling {
					for _, a1 := range e.Attr {
						if a1.Key == "title"{
							fmt.Println("time: ", a1.Val)
						}
						if a1.Key == "href"{
							fmt.Println("link: ", a1.Val)
						}
					}

				}

			}
		}
	}
}
