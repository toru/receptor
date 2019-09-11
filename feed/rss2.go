package feed

import (
	"encoding/xml"
	"time"

	"github.com/toru/dexter/index"
)

// RSS2Time is a custom type that embeds the standard time.Time.
// Purpose of this type is to implement a custom XML node parser.
type RSS2Time struct {
	time.Time
}

// UnmarshalXML implements the xml.Unmarshaler interface.
// A generous implementation that allows multiple format strings.
func (rt *RSS2Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	var tmp string
	var parsed time.Time

	d.DecodeElement(&tmp, &start)
	fmts := []string{time.RFC1123Z, time.RFC1123, time.RFC822Z, time.RFC822}
	for _, f := range fmts {
		parsed, err = time.Parse(f, tmp)
		if err == nil {
			rt.Time = parsed
			break
		}
	}
	return err
}

type RSS2Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	GUID        string `xml:"guid"`

	// Dexter specific
	feedID index.ID
}

type RSS2Channel struct {
	// Required Fields
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`

	// Optional Fields
	Language       string   `xml:"language"`
	Copyright      string   `xml:"copyright"`
	ManagingEditor string   `xml:"managingEditor"`
	WebMaster      string   `xml:"webMaster"`
	PubDate        RSS2Time `xml:"pubDate"`
	LastBuildDate  RSS2Time `xml:"lastBuildDate"`
	Category       []string `xml:"category"`
	Generator      string   `xml:"generator"`
	Docs           string   `xml:"docs"`
	Cloud          string   `xml:"cloud"`
	TTL            uint     `xml:"ttl"`
	Image          string   `xml:"image"`
	// TODO: TextInput
	// TODO: SkipHours
	// TODO: SkipDays

	// Dexter specific
	Items []RSS2Item `xml:"item"`
}

type RSS2Feed struct {
	XMLName xml.Name    `xml:"rss"`
	Version string      `xml:"version,attr"`
	Channel RSS2Channel `xml:"channel"`

	// Dexter specific attributes
	subscriptionID index.DexID
}

// ParseRSS2Feed parses the given byte slice as an RSS2Feed.
func ParseRSS2Feed(doc []byte) (Feed, error) {
	feed := &RSS2Feed{}
	if err := xml.Unmarshal(doc, feed); err != nil {
		return nil, err
	}
	return feed, nil
}

// ID implements the Feed interface. RSS 2.0 doesn't define a feed
// identifier so return the subscription_id as a hex string instead.
func (rf *RSS2Feed) ID() string {
	if len(rf.subscriptionID) > 0 {
		return index.DexIDToHexDigest(rf.subscriptionID[:])
	}
	return ""
}

// Title implements the Feed interface.
func (rf *RSS2Feed) Title() string {
	return rf.Channel.Title
}

// Format implements the Feed interface.
func (rf *RSS2Feed) Format() uint {
	return RSS2FeedFormat
}

// SubscriptionID implements the Feed interface.
func (rf *RSS2Feed) SubscriptionID() []byte {
	return rf.subscriptionID[:]
}

// Entries implements the Feed interface.
func (rf *RSS2Feed) Entries() []Entry {
	rv := make([]Entry, len(rf.Channel.Items))
	for i := range rf.Channel.Items {
		rv[i] = &rf.Channel.Items[i]
	}
	return rv
}

// SetSubscriptionID sets the given ID to the feed.
func (rf *RSS2Feed) SetSubscriptionID(id []byte) {
	var tmp [index.SHA224DexIDLen]byte
	copy(tmp[:], id)
	rf.subscriptionID = tmp
}

// SetFeedID implements the Entry interface.
func (ri *RSS2Item) SetFeedID(id index.ID) {
	ri.feedID = id
}

// FeedID implements the Entry interface.
func (ri *RSS2Item) FeedID() index.ID {
	return ri.feedID
}

// ID implements the Entry interface.
func (ri *RSS2Item) ID() string {
	return ri.GUID
}

// Summary implements the Entry interface.
func (ri *RSS2Item) Summary() string {
	return ri.Description
}
