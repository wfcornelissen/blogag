package rss

import "fmt"

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (rf *RSSFeed) Display() {
	fmt.Printf("Channel Title: %v\n", rf.Channel.Title)
	fmt.Printf("Channel Link: %v\n", rf.Channel.Link)
	fmt.Printf("Channel Description: %v\n", rf.Channel.Description)
	fmt.Println("____________________________________________________")
	fmt.Println("____________________________________________________")
	if len(rf.Channel.Item) > 0 {
		for _, item := range rf.Channel.Item {
			fmt.Printf("Item Title: %v\n", item.Title)
			fmt.Printf("Item Link: %v\n", item.Link)
			fmt.Printf("Item Description: %v\n", item.Description)
			fmt.Printf("Item Publish Date: %v\n", item.PubDate)
			fmt.Println("____________________________________________________")

		}
	}
	fmt.Println("-----END OF RSS FEED-----")

}
