package feedly

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type FeedContent struct {
	Title string     `json:"title"`
	Items []FeedItem `json:"items"`
}

type FeedItem struct {
	Title    string `json:"title"`
	OriginId string `json:"originId"`
}

func StreamContent(streamId string, ch chan<- FeedContent) {
	// TODO: proper http resp code handling
	url := "https://cloud.feedly.com/v3/streams/contents?streamId=" + streamId
	fmt.Printf("Requesting %s \n", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var feed FeedContent
	err = json.Unmarshal(respBytes, &feed)

	fmt.Printf("Response %v \n", feed)

	if err != nil {
		panic(err)
	}
	ch <- feed
}

func StreamContents(streamIds []string, ch chan<- FeedContent) {
	for _, id := range streamIds {
		go StreamContent(id, ch)
	}
}
