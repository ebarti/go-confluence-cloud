package confluentcloud

import (
	"net/http"
	"net/url"
)

type API interface {
	Request(*http.Request) ([]byte, error)
	SendContentRequest(*url.URL, string, *Content) (*Content, error)
	VerifyTLS(bool)
	GetContent(ContentQuery) (*Content, error)
	GetContentFromNext(Links) (*Content, error)
	GetAttachmentsFromResult(Results, string) ([]Results, error)
}

// api is the main api data structure
type api struct {
	endPoint        *url.URL
	client          *http.Client
	username, token string
}

type Content struct {
	Results []Results `json:"results,omitempty"`
	Start   int       `json:"start,omitempty"`
	Limit   int       `json:"limit,omitempty"`
	Size    int       `json:"size,omitempty"`
	Links   Links     `json:"_links,omitempty"`
}

type Children struct {
	Attachment Attachment `json:"attachment,omitempty"`
}

type Attachment struct {
	Results []Results `json:"results,omitempty"`
	Start   int       `json:"start,omitempty"`
	Limit   int       `json:"limit,omitempty"`
	Size    int       `json:"size,omitempty"`
	Links   Links     `json:"_links,omitempty"`
}

type Results struct {
	ID         string     `json:"id,omitempty"`
	Type       string     `json:"type,omitempty"`
	Status     string     `json:"status,omitempty"`
	Title      string     `json:"title,omitempty"`
	Children   Children   `json:"children,omitempty"`
	Body       Body       `json:"body,omitempty"`
	Expandable Expandable `json:"_expandable,omitempty"`
	Links      Links      `json:"_links,omitempty"`
	Metadata   Metadata   `json:"metadata,omitempty"`
}

type Storage struct {
	Value           string        `json:"value,omitempty"`
	Representation  string        `json:"representation,omitempty"`
	Embeddedcontent []interface{} `json:"embeddedContent,omitempty"`
	Expandable      Expandable    `json:"_expandable,omitempty"`
}

type Body struct {
	Storage Storage `json:"storage,omitempty"`
}

type Metadata struct {
	MediaType string `json:"mediaType,omitempty"`
}

type Expandable struct {
	Space string `json:"space,omitempty"`
}

type Links struct {
	Base     string `json:"base,omitempty"`
	Self     string `json:"self,omitempty"`
	Next     string `json:"next,omitempty"`
	Tinyui   string `json:"tinyui,omitempty"`
	Webui    string `json:"webui,omitempty"`
	Download string `json:"download,omitempty"`
}

// ContentQuery defines the query parameters
// used for content related searching
// Query parameter values https://developer.atlassian.com/cloud/confluence/rest/#api-content-get
type ContentQuery struct {
	Expand     []string
	Limit      int    // page limit
	OrderBy    string // fieldpath asc/desc e.g: "history.createdDate desc"
	PostingDay string // required for blogpost type Format: yyyy-mm-dd
	SpaceKey   string
	Start      int    // page start
	Status     string // current, trashed, draft, any
	Title      string // required for page
	Trigger    string // viewed
	Type       string // page, blogpost
	Version    int    //version number when not lastest
}
