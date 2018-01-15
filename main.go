package main

import (
	"github.com/jezzay/xkcd/feedly"
	"fmt"
)

func main() {

	contents := feedly.StreamContents("feed/http://xkcd.com/rss.xml")
	fmt.Printf("%+v", contents)

}
