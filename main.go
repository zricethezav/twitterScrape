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
       hashTags []string
       tweetUrl string
       time float64 // unix time
}

func main() {
       // this program demonstrates scraping
       resp, err := http.Get("https://twitter.com/realDonaldTrump")
       if err != nil {
              panic(err)
       }
       root, err := html.Parse(resp.Body)
       if err != nil {
              panic(err)
       }

       /*
       var w func(*html.Node)
       w = func(n *html.Node) {
              if n.Type == html.ElementNode && n.Data == "div" {
                     for _, a := range n.Attr {
                            if a.Key == "class" {
                                   if a.Val == "stream-item-header" {
                                          fmt.Println("YOOOOOOOOOOOOO FUCK")
                                          fmt.Println(a.Val)
                                          break
                                   }
                            }
                     }
              }
       }
       */

       var f func(*html.Node, bool)
       f = func(n *html.Node, parentTweet bool)  {
	       /*
              if n.Type == html.TextNode && parentTweet {
                     fmt.Println(n.Data)
                     for _, a := range n.Attr {
                            fmt.Println(a.Key)
                     }
              }
              */

              if n.Type == html.ElementNode && n.Data == "div" {
                     for _, a := range n.Attr {
                            if a.Key == "class" {
                                   if a.Val == "js-tweet-text-container" {
					   // process tweet
					   fmt.Println(p(n, ""))
					   parentTweet = true
					   break
                                   }
                            }
                     }
              }

              for c := n.FirstChild; c != nil; c = c.NextSibling {
                     f(c, parentTweet)
              }
       }

       f(root, false)
}

func p(n *html.Node, tweet string) string {
	if n.Type == html.TextNode {
		fmt.Println(n.Data)
		for _, a := range n.Attr {
			fmt.Println(a.Key)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tweet += p(c, "hhehe")
	}
	return tweet
}
