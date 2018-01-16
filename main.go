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
	feedIds := flag.String("feedIds", "feed/http://xkcd.com/rss.xml", "Feed ids to search on")

	flag.Parse()

	if *recent && len(*feedIds) >= 1 {
		searchFeeds(*feedIds)
	} else {
		fmt.Printf("-feedIds needs to be provided")
		os.Exit(-1)
	}

}

func searchFeeds(feedIds string) {
	feeds := make(chan feedly.FeedContent)
	done := make(chan bool, 2)

	feedLen := 1
	// TODO: split feed ids on comma

	// Issue multiple requests and wait
	feedly.StreamContents([]string{feedIds}, feeds)

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
