package confluentcloud

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func Test_addSearchQueryParams(t *testing.T) {
	type args struct {
		query SearchContentQuery
	}
	tests := []struct {
		name string
		args args
		want *url.Values
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addSearchQueryParams(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addSearchQueryParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_api_GetSearchContentResults(t *testing.T) {
	server := confluenceRestAPIStub()
	defer server.Close()

	api, err := NewAPI(server.URL+"/wiki/rest/api", "username", "token")
	assert.Nil(t, err)

	s, err := api.GetSearchContentResults(SearchContentQuery{})
	assert.Nil(t, err)
}

func Test_api_getSearchEndpoint(t *testing.T) {
	type fields struct {
		endPoint *url.URL
		client   *http.Client
		username string
		token    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *url.URL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &api{
				endPoint: tt.fields.endPoint,
				client:   tt.fields.client,
				username: tt.fields.username,
				token:    tt.fields.token,
			}
			got, err := a.getSearchEndpoint()
			if (err != nil) != tt.wantErr {
				t.Errorf("getSearchEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSearchEndpoint() got = %v, want %v", got, tt.want)
			}
		})
	}
}
