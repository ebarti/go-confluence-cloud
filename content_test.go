package confluentcloud

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetContentEndpoint(t *testing.T) {
	a, err := newAPI("https://test.test", "username", "token")
	assert.Nil(t, err)

	url, err := a.getContentEndpoint()
	assert.Nil(t, err)
	assert.Equal(t, "/content/", url.Path)
}

func Test_GetContentIDEndpoint(t *testing.T) {
	a, err := newAPI("https://test.test", "username", "token")
	assert.Nil(t, err)

	url, err := a.getContentIDEndpoint("test")
	assert.Nil(t, err)
	assert.Equal(t, "/content/test", url.Path)
}

func Test_GetContentChildEndpoint(t *testing.T) {
	a, err := newAPI("https://test.test", "username", "token")
	assert.Nil(t, err)

	url, err := a.getContentChildEndpoint("1", "2")
	assert.Nil(t, err)
	assert.Equal(t, "/content/1/child/2", url.Path)
}

func Test_GetContentGenericEndpoint(t *testing.T) {
	a, err := newAPI("https://test.test", "username", "token")
	assert.Nil(t, err)

	url, err := a.getContentGenericEndpoint("1", "2")
	assert.Nil(t, err)
	assert.Equal(t, "/content/1/2", url.Path)
}

func Test_GetContent(t *testing.T) {
	server := confluenceRestAPIStub()
	defer server.Close()

	api, err := NewAPI(server.URL+"/wiki/rest/api", "username", "token")
	assert.Nil(t, err)

	s, err := api.GetContent(ContentQuery{})
	assert.Nil(t, err)
	assert.Equal(t, &Content{
		Results: []Results{{
			ID: "ContentResult",
			Children: Children{Attachment: Attachment{
				Results: []Results{AttachmentResult},
				Links:   Links{Next: "/rest/api/content/1/child/attachment?limit=25&start=25"},
			}},
		}},
		Links: Links{Base: "http://" + URL + "/wiki"},
	}, s)
}

func Test_GetContenWithClient(t *testing.T) {
	server := confluenceRestAPIStub()
	defer server.Close()
	myClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}}
	api, err := NewAPIWithClient(server.URL+"/wiki/rest/api", myClient)
	assert.Nil(t, err)

	s, err := api.GetContent(ContentQuery{})
	assert.Nil(t, err)
	assert.Equal(t, &Content{
		Results: []Results{{
			ID: "ContentResult",
			Children: Children{Attachment: Attachment{
				Results: []Results{AttachmentResult},
				Links:   Links{Next: "/rest/api/content/1/child/attachment?limit=25&start=25"},
			}},
		}},
		Links: Links{Base: "http://" + URL + "/wiki"},
	}, s)
}

func Test_GetContentFromNext(t *testing.T) {
	server := confluenceRestAPIStub()
	defer server.Close()

	api, err := NewAPI(server.URL+"/wiki/rest/api", "username", "token")
	assert.Nil(t, err)

	links := Links{Base: server.URL + "/wiki", Next: "/rest/api/content?limit=25&start=0"}

	s, err := api.GetContentFromNext(links)
	assert.Nil(t, err)
	assert.Equal(t, &Content{
		Results: []Results{{
			ID: "ContentResult",
			Children: Children{Attachment: Attachment{
				Results: []Results{AttachmentResult},
				Links:   Links{Next: "/rest/api/content/1/child/attachment?limit=25&start=25"},
			}},
		}},
		Links: Links{Base: "http://" + URL + "/wiki"},
	}, s)
}

func TestAttachmentGetter(t *testing.T) {
	server := confluenceRestAPIStub()
	defer server.Close()

	api, err := NewAPI(server.URL+"/wiki/rest/api", "username", "token")
	assert.Nil(t, err)

	s, err := api.GetContent(ContentQuery{})
	assert.Nil(t, err)
	assert.Equal(t, "/rest/api/content/1/child/attachment?limit=25&start=25", s.Results[0].Children.Attachment.Links.Next)
	r, err := api.GetAttachmentsFromResult(s.Results[0], s.Links.Base)
	assert.Equal(t, 3, len(r))
}

func TestAddContentQueryParams(t *testing.T) {
	query := ContentQuery{
		Expand:     []string{"foo", "bar"},
		Limit:      1,
		OrderBy:    "test",
		PostingDay: "test",
		SpaceKey:   "test",
		Start:      1,
		Status:     "test",
		Title:      "test",
		Trigger:    "test",
		Type:       "test",
	}

	p := addContentQueryParams(query)

	assert.Equal(t, p.Get("expand"), "foo,bar")
	assert.Equal(t, p.Get("limit"), "1")
	assert.Equal(t, p.Get("orderby"), "test")
	assert.Equal(t, p.Get("postingDay"), "test")
	assert.Equal(t, p.Get("spaceKey"), "test")
	assert.Equal(t, p.Get("start"), "1")
	assert.Equal(t, p.Get("status"), "test")
	assert.Equal(t, p.Get("title"), "test")
	assert.Equal(t, p.Get("trigger"), "test")
	assert.Equal(t, p.Get("type"), "test")
}
