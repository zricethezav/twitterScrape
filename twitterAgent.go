package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"time"
	"strconv"
	"strings"
	"regexp"
)

// constants
const (
	twitterTimeStringFmt = "15:04 PM - 2 Jan 2006"
	localTimeZone        = "CST" // timezone where the compute is happening
	hashTagRegEx = "#([a-zA-Z\\d]+)"
	mentionRegEx = "@([a-zA-Z\\d]+)"
	twitterLoadMoreBase = "https://twitter.com/i/profiles/show/%s/timeline/tweets?include_available_features=1&include_entities=1&max_position=%s&reset_error_state=false"
)

type TweetHeader struct {
	tweetUrl   string
	time       int64 // unix time
	stringTime string
}
type TweetFooter struct {
	favorites int
	retweets  int
	replies   int
}
type Tweet struct {
	tweetText   string
	hashTags    []string
	mentions []string
	tweetHeader TweetHeader
	tweetFooter TweetFooter
}

func getTweets(twitterHandle string ) []Tweet {
	resp, err := http.Get(fmt.Sprintf("https://twitter.com/%s", twitterHandle))
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

	// htmlTweetParser traverses html tags which is a tree structure. It looks for the tweet contents, header, and
	// footer and will call their respective helper functions to retrieve data for the content, header, and footer.
	var htmlTweetParser func(*html.Node)
	htmlTweetParser = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" {
					if a.Val == "stream-item-header" {
						currTweet.tweetHeader = retrieveTweetHeader(n.FirstChild)
						currEleTweet = true
					} else if a.Val == "js-tweet-text-container" {
						currTweet.tweetText = retrieveTweet(n)
						currTweet.hashTags = extractHashTags(currTweet.tweetText)
						currTweet.mentions = extractMentions(currTweet.tweetText)
					} else if a.Val == "stream-item-footer" {
						currTweet.tweetFooter = retrieveTweetFooter(n.FirstChild)
					}
				}
			}
			// if current node contains attributes that signal a tweet, add to
			// tweet slice and reset currEleTweet flag
			if currEleTweet {
				currEleTweet = false
				tweets = append(tweets, currTweet)
			}
		}

		// DFS
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlTweetParser(c)
		}
	}
	htmlTweetParser(root)
	return tweets
}

// extractHashTags
func extractHashTags(tweetText string) []string {
	var hashTags []string
	if !strings.Contains(tweetText, "#") {
		return hashTags
	} else {
		// regex matching
		// https://regex-golang.appspot.com/assets/html/index.html
		r, _ := regexp.Compile(hashTagRegEx)
		hashTags = r.FindAllString(tweetText, -1)
		return hashTags
	}
}

//extractMentions
func extractMentions(tweetText string) []string {
	var mentions []string
	if !strings.Contains(tweetText, "@") {
		return mentions
	} else {
		r, _ := regexp.Compile(mentionRegEx)
		mentions = r.FindAllString(tweetText, -1)
		return mentions
	}
}

// retrieveTweet looks for all TextNodes in the js-tweet-text-container tag and constructs the tweet.
// returns a String
func retrieveTweet(n *html.Node) string {
	tweet := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tweet += retrieveTweet(c)
	}
	if n.Type == html.TextNode {
		partialTweet := string(n.Data)
		// preventing run-on hashtags/links
		if strings.Contains(partialTweet, "#") || strings.Contains(partialTweet, "@"){
			return partialTweet
		} else {
			return partialTweet + " "
		}
	}
	return tweet
}

// retrieveTweetHeader parses the children of an html.Node object and extracts the timestamp and url of the tweet.
// Returns a TweetHeader
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

	for c := n; c != nil; c = c.NextSibling {
		for _, a := range c.Attr {
			if a.Val == "time" {
				processTweetHeaderHelper(c)
			}
		}
	}
	return header
}

// retrieveTweetFooter parses the children of an html.Node object and extracts the number of likes, replies,
// and retweets.
// Returns a TweetFooter
func retrieveTweetFooter(n *html.Node) TweetFooter {
	var footer TweetFooter

	// processTweetFooterHelper is the ProfileTweet-actionCountlist processor.
	// It drills down and looks for replies, likes, and retweets.
	var processTweetFooterHelper func(n *html.Node)
	processTweetFooterHelper = func(n *html.Node) {
		for e := n.FirstChild; e != nil; e = e.NextSibling {
			for _, a1 := range e.Attr {
				if a1.Val == "ProfileTweet-action--reply u-hiddenVisually" {
					count, err := strconv.Atoi(e.FirstChild.NextSibling.Attr[1].Val)
					if err != nil {
						footer.replies = 0
					} else {
						footer.replies = count
					}
				}
				if a1.Val == "ProfileTweet-action--favorite u-hiddenVisually" {
					count, err := strconv.Atoi(e.FirstChild.NextSibling.Attr[1].Val)
					if err != nil {
						footer.favorites = 0
					} else {
						footer.favorites = count
					}
				}
				if a1.Val == "ProfileTweet-action--retweet u-hiddenVisually" {
					count, err := strconv.Atoi(e.FirstChild.NextSibling.Attr[1].Val)
					if err != nil {
						footer.retweets = 0
					} else {
						footer.retweets = count
					}
				}
			}
		}
	}
	for c := n; c != nil; c = c.NextSibling {
		for _, a := range c.Attr {
			if a.Val == "ProfileTweet-actionCountList u-hiddenVisually" {
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
