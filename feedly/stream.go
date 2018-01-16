package feedly

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type FeedContent struct {
	Title string     `json:"title"`
	Items []FeedItem `json:"items"`
}

type FeedItem struct {
	Title    string `json:"title"`
	OriginId string `json:"originId"`
}

func StreamContent(streamId string, apiToken string, ch chan<- FeedContent) {
	// TODO: proper http resp code handling

	url := "https://cloud.feedly.com/v3/streams/contents?streamId=" + streamId

	client := &http.Client{Timeout: time.Second * 5}

	req, _ := http.NewRequest("GET", url, nil)
	if len(apiToken) >= 1 {
		req.Header.Add("Authorization", "Bearer "+apiToken)
	}

	//fmt.Printf("Requesting %s \n", url)
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	//fmt.Printf("Full API Response %s \n", respBytes)

	defer resp.Body.Close()

	var feed FeedContent
	err = json.Unmarshal(respBytes, &feed)

	// fmt.Printf("Response %v \n", feed)

	if err != nil {
		panic(err)
	}
	ch <- feed
}

func StreamContents(streamIds []string, apiToken string, ch chan<- FeedContent) {
	for _, id := range streamIds {
		go StreamContent(id, apiToken, ch)
	}
}
