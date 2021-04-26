package confluentcloud

import (
	"net/http"
	"net/url"
	"time"
)

type API interface {
	Request(*http.Request) ([]byte, error)
	SendContentRequest(*url.URL, string, *Content) (*Content, error)
	VerifyTLS(bool)
	GetContent(ContentQuery) (*Content, error)
	GetContentFromNext(Links) (*Content, error)
	GetAttachmentsFromResult(Results, string) ([]Results, error)
	GetSearchContentResults(SearchContentQuery) (*SearchPageResults, error)
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

type SearchContentQuery struct {
	Cql                   string
	CqlContext            map[string]string // The space, content, and content status to execute the search against: spaceKey, contentId, contentStatuses
	cursor                string            // Pointer to a set of search results, returned as part of the next or prev URL from the previous search call.
	Limit                 int               // The maximum number of content objects to return per page. Note, this may be restricted by fixed system limits.
	IncludeArchivedSpaces bool              // Include content from archived spaces in the results.
}

type SearchPageResults struct {
	Results        []SearchPageResult `json:"results"`
	Start          int                `json:"start"`
	Limit          int                `json:"limit"`
	Size           int                `json:"size"`
	TotalSize      int                `json:"totalSize"`
	CqlQuery       string             `json:"cqlQuery"`
	SearchDuration int                `json:"searchDuration"`
	Links          Links              `json:"_links"`
}

type SearchPageResult struct {
	Content               Content          `json:"content"`
	Title                 string           `json:"title"`
	Excerpt               string           `json:"excerpt"`
	Url                   string           `json:"url"`
	ResultParentContainer ContainerSummary `json:"resultParentContainer"`
	ResultGlobalContainer ContainerSummary `json:"resultGlobalContainer"`
	Breadcrumbs           []Breadcrumb     `json:"breadcrumbs"`
	EntityType            string           `json:"entityType"`
	IconCssClass          string           `json:"iconCssClass"`
	LastModified          time.Time        `json:"lastModified"`
	FriendlyLastModified  string           `json:"friendlyLastModified"`
}
type ContainerSummary struct {
	Title      string `json:"title"`
	DisplayUrl string `json:"displayUrl"`
}

// Breadcrumb struct for Breadcrumb
type Breadcrumb struct {
	Label     string `json:"label"`
	Url       string `json:"url"`
	Separator string `json:"separator"`
}
