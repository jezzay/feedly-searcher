package main

import (
	"github.com/jezzay/feedly-searcher/feedly"
	"flag"
	"fmt"
	"os"
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

	if *recent && len(*feedIds) >= 1 {
		searchFeeds(*feedIds, apiToken)
	} else {
		fmt.Printf("-feedIds needs to be provided")
		os.Exit(-1)
	}

}

func searchFeeds(feedIds string, apiToken string) {
	feeds := make(chan feedly.FeedContent)
	done := make(chan bool, 2)

	feedLen := 1
	// TODO: split feed ids on comma

	// Issue multiple requests and wait
	feedly.StreamContents([]string{feedIds},apiToken, feeds)

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
			if received == feedLen {
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
