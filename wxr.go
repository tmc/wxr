// Package wxr implements types corresponding to the Wordpress WXR XML encoding.
// The initial implementation is consumption oriented and may be incomplete.
package wxr

import (
	"encoding/xml"
	"io"
	"time"
)

// RSS is the root element of an WXR document
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel
}

// Category describes the categories in the blog export
type Category struct {
	ID   int    `xml:"term_id"`
	Slug string `xml:"category_nicename"`
	Name string `xml:"cat_name"`
}

// Channel is the element describing the blog
type Channel struct {
	XMLName     xml.Name   `xml:"channel"`
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Categories  []Category `xml:"category"`
	Description string     `xml:"description"`
	WXRVersion  string     `xml:"wxr_version"`
	Items       []Item     `xml:"item"`
}

// Item is a blog post
type Item struct {
	XMLName    xml.Name       `xml:"item"`
	ID         int            `xml:"post_id"`
	Name       string         `xml:"post_name"`
	Title      string         `xml:"title"`
	Link       string         `xml:"link"`
	Categories []ItemCategory `xml:"category"`
	Type       string         `xml:"post_type"`
	PubDate    WpTime         `xml:"pubDate"`
}

// ItemCategory is a Category that is associated with an Item
type ItemCategory struct {
	XMLName xml.Name `xml:"category"`
	Slug    string   `xml:"nicename,attr"`
	Name    string   `xml:",chardata"`
}

// WpTime exists to provide UnMarshaling for the wordpress pubDate format
type WpTime time.Time

// UnmarshalText attempts to unmarshall the provided byte slice as a time.RFC1123Z
func (t *WpTime) UnmarshalText(text []byte) error {
	// Mon Jan 2 15:04:05 -0700 MST 2006
	// Mon, 03 Sep 2007 18:23:34 +0000
	parsed, err := time.Parse(time.RFC1123Z, string(text))
	*t = WpTime(parsed)
	return err
}

func (t WpTime) String() string {
	return time.Time(t).String()
}

// NewRSS attempts to parse the provided Reader into an RSS instance
func NewRSS(r io.Reader) (*RSS, error) {
	result := new(RSS)
	d := xml.NewDecoder(r)
	err := d.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
