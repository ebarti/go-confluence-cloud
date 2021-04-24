package confluentcloud

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"reflect"
	"testing"
)

func Test_api_GetSearchContentResults(t *testing.T) {
	server := confluenceRestAPIStub()
	defer server.Close()

	api, err := NewAPI(server.URL+"/wiki/rest/api", "username", "token")
	assert.Nil(t, err)

	s, err := api.GetSearchContentResults(SearchContentQuery{})
	assert.Nil(t, err)
	assert.Equal(t, &SearchPageResults{
		Results:        []SearchPageResult{},
		TotalSize:      1200,
		CqlQuery:       "testQUERY",
		SearchDuration: 100,
		Links: Links{
			Base: "base",
		},
	}, s)
}

func Test_api_getSearchEndpoint(t *testing.T) {
	api, err := newAPI("https://test.test", "username", "token")
	assert.Nil(t, err)
	url, err := api.getSearchEndpoint()
	assert.Nil(t, err)
	assert.Equal(t, "/search", url.Path)
}

func Test_addSearchQueryParams(t *testing.T) {
	cql := "type IN (blogpost, page)"
	query := SearchContentQuery{
		Cql:   cql,
		Limit: 1,
	}
	got := addSearchQueryParams(query)
	want := url.Values{}
	want.Set("cql", cql)
	want.Set("limit", "1")
	if !reflect.DeepEqual(got, &want) {
		t.Errorf("addSearchQueryParams() = %v, want %v", got, &want)
	}
}
