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

type NewRequest struct {
	StreamIds []string
	ApiToken  string
	Unread    bool
}

func StreamContent(feedReq NewRequest, streamId string, ch chan<- FeedContent) {
	// TODO: proper http resp code handling

	url := "https://cloud.feedly.com/v3/streams/contents?streamId=" + streamId
	if feedReq.Unread {
		url += "&unreadOnly=true"
	}

	client := &http.Client{Timeout: time.Second * 5}
	httpReq, _ := http.NewRequest("GET", url, nil)

	if len(feedReq.ApiToken) >= 1 {
		httpReq.Header.Add("Authorization", "Bearer "+feedReq.ApiToken)
	}

	//fmt.Printf("Requesting %s \n", url)
	resp, err := client.Do(httpReq)

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

func StreamContents(req NewRequest, ch chan<- FeedContent) {
	for _, id := range req.StreamIds {
		go StreamContent(req, id, ch)
	}
}
