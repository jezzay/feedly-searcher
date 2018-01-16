package main

import (
	"github.com/jezzay/feedly-searcher/feedly"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	recent := flag.Bool("recent", true, "List most recent posts")
	unreadOnly := flag.Bool("unread", false,
		`List only the posts have not been marked as read in Feedly. Requires a feedly oAuth token to be set in an env var called feedly-api-token`)
	feedIds := flag.String("feedIds", "feed/http://xkcd.com/rss.xml", "Feed ids to search on")

	flag.Parse()

	apiToken := oAuthToken()

	if *unreadOnly && len(apiToken) == 0 {
		fmt.Printf("An feedly oAuth token needs to be set in a env var called feedly-api-token in order to use the -unread option")
		os.Exit(-1)
	}

	req := feedly.NewRequest{StreamIds: strings.Split(*feedIds, ","), ApiToken: apiToken, Unread: *unreadOnly}

	if *recent && len(*feedIds) >= 1 {
		searchFeeds(req)
	} else {
		fmt.Printf("-feedIds needs to be provided")
		os.Exit(-1)
	}

}

func searchFeeds(req feedly.NewRequest) {
	feeds := make(chan feedly.FeedContent)
	done := make(chan bool, 1)

	feedReqLen := len(req.StreamIds)

	// Issue multiple requests and wait
	feedly.StreamContents(req, feeds)

	received := 0
	// non blocking select
	for {
		select {
		case contents := <-feeds:
			fmt.Printf("Most recent posts for %s: \n", contents.Title)

			for _, f := range contents.Items {
				fmt.Printf("%s available at %s\n", f.Title, f.OriginId)
			}
			//fmt.Printf("\n Result number %v from multiple requests:  %+v \n", received, contents)
			received++
			if received == feedReqLen {
				done <- true
			}
		case <-done:
			// fmt.Printf("All channels have completed; exiting")
			return
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func oAuthToken() string {
	return os.Getenv("feedly-api-token")
}
