package feedly

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type FeedContent struct {
	Title string `json:"title"`
	Items []FeedItem `json:"items"`
}

type FeedItem struct {
	Title string `json:"title"`
	OriginId string `json:"originId"`
}

func StreamContents(streamId string) FeedContent  {
	// TODO: proper http resp code handling
	resp, err := http.Get("https://cloud.feedly.com/v3/streams/contents?streamId=" + streamId)
	if err != nil {
		panic(err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var feed FeedContent
	err = json.Unmarshal(respBytes, &feed)
	if err != nil {
		panic(err)
	}
	return feed
}