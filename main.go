package main

import (
	"github.com/jezzay/xkcd/feedly"
	"fmt"
	"time"
)

func main() {
	feed := make(chan feedly.FeedContent)
	feeds := make(chan feedly.FeedContent)
	done := make(chan bool, 2)

	// Issue multiple requests and wait
	feedly.StreamContents([]string{"feed/http://xkcd.com/rss.xml", "feed/http://xkcd.com/rss.xml"}, feeds)

	// Single request
	go feedly.StreamContent("feed/http://xkcd.com/rss.xml", feed)

	received := 0
	completedChannels := 0

	// non blocking select
	for {
		select {
		case content := <-feed:
			fmt.Printf("\n Result from single request =  %+v \n", content)
			done <- true
		case contents := <-feeds:
			fmt.Printf("\n Result number %v from multiple requests:  %+v \n", received, contents)
			received++
			if received == 2 {
				done <- true
			}
		case <-done:
			completedChannels++
			if completedChannels == 2 {
			fmt.Printf("All channels have completed; exiting")
				return
			}

		default:
			time.Sleep(50 * time.Millisecond)
		}
	}

	// blocking read from channels; both the single and range will block

	//content := <-feed
	//fmt.Printf("\n Result from single request =  %+v \n", content)
	//
	//received := 0
	//for content := range feeds {
	//	fmt.Printf("\n Result number %v from multiple requests:  %+v \n", received, content)
	//	received++
	//	if received == 2 {
	//		close(feeds)
	//	}
	//}
}
