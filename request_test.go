package confluentcloud

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testValuesRequest struct {
	Endpoint string
	Body     string
	Error    error
}

var AttachmentResult Results = Results{
	ID: "AttachmentResult",
}

const URL string = "127.0.0.1:52347"

func TestRequest(t *testing.T) {
	server := confluenceRestAPIStub()
	defer server.Close()

	api, err := newAPI(server.URL+"/wiki/rest/api", "userame", "token")
	assert.Nil(t, err)

	testValues := []testValuesRequest{
		{"/test", "\"test\"", nil},
		{"/nocontent", "", nil},
		{"/noauth", "", fmt.Errorf("authentication failed")},
		{"/noservice", "", fmt.Errorf("service is not available: 503 Service Unavailable")},
		{"/internalerror", "", fmt.Errorf("internal server error: 500 Internal Server Error")},
		{"/unknown", "", fmt.Errorf("unknown response status: 408 Request Timeout")},
	}

	for _, test := range testValues {

		req, err := http.NewRequest(http.MethodGet, api.endPoint.String()+test.Endpoint, nil)
		assert.Nil(t, err)

		b, err := api.Request(req)
		if test.Error == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.Error.Error(), err.Error())
		}

		assert.Equal(t, string(b), test.Body)
	}
}

func confluenceRestAPIStub() *httptest.Server {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp interface{}
		switch r.RequestURI {
		case "/wiki/rest/api/test":
			resp = "test"
		case "/wiki/rest/api/noauth":
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		case "/wiki/rest/api/nocontent":
			http.Error(w, "", http.StatusNoContent)
			return
		case "/wiki/rest/api/noservice":
			http.Error(w, "", http.StatusServiceUnavailable)
			return
		case "/wiki/rest/api/internalerror":
			http.Error(w, "", http.StatusInternalServerError)
			return
		case "/wiki/rest/api/unknown":
			http.Error(w, "", http.StatusRequestTimeout)
			return
		case "/wiki/rest/api/search":
			resp = SearchPageResults{
				Results:        []SearchPageResult{},
				TotalSize:      1200,
				CqlQuery:       "testQUERY",
				SearchDuration: 100,
				Links: Links{
					Base: "base",
				},
			}
		case "/wiki/rest/api/content/", "/wiki/rest/api/content?limit=25&start=0":
			resp = Content{
				Results: []Results{{
					ID: "ContentResult",
					Children: Children{Attachment: Attachment{
						Results: []Results{AttachmentResult},
						Links:   Links{Next: "/rest/api/content/1/child/attachment?limit=25&start=25"},
					}},
				}},
				Links: Links{Base: "http://" + URL + "/wiki"},
			}
		case "/wiki/rest/api/content?next=true&expand=body.storage&limit=1&start=1":
			resp = Content{
				Results: []Results{{
					ID: "ContentResult",
					Children: Children{Attachment: Attachment{
						Results: []Results{AttachmentResult},
						Links:   Links{Next: "/rest/api/content/1/child/attachment?limit=25&start=25"},
					}},
				}},
				Links: Links{Base: "http://" + URL + "/wiki"},
			}
		case "/wiki/rest/api/content/1/child/attachment?limit=25&start=25":
			resp = Content{
				Results: []Results{AttachmentResult, AttachmentResult},
			}
		default:
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}))
	l, err := net.Listen("tcp", URL)
	if err != nil {
		return nil
	}
	server.Listener.Close()
	server.Listener = l
	server.Start()
	return server
}

func Test_SendContentRequest(t *testing.T) {

	server := confluenceRestAPIStub()
	defer server.Close()

	api, err := newAPI(server.URL+"/wiki/rest/api", "userame", "token")
	_ = server.Config.Addr
	assert.Nil(t, err)

	type args struct {
		c      *Content
		path   string
		method string
	}
	tests := []struct {
		name    string
		args    args
		want    *Content
		wantErr bool
	}{
		{
			name: "IfCantReachPath_ReturnsError",
			args: args{
				path:   "786/wiki/rest/api/content/",
				method: http.MethodGet,
			},
			wantErr: true,
		},
		{
			name: "IfCantReachPath_ReturnsError",
			args: args{
				c:      &Content{},
				path:   "786/wiki/rest/api/content/",
				method: http.MethodGet,
			},
			wantErr: true,
		},
		{
			name: "IfServerDoesntReturnContent_ReturnsError",
			args: args{
				path:   "/wiki/rest/api/test",
				method: http.MethodGet,
			},
			wantErr: true,
		},
		{
			name: "When GET - returns correct response",
			args: args{
				path:   "/wiki/rest/api/content/",
				method: http.MethodGet,
			},
			want: &Content{
				Results: []Results{{
					ID: "ContentResult",
					Children: Children{Attachment: Attachment{
						Results: []Results{AttachmentResult},
						Links:   Links{Next: "/rest/api/content/1/child/attachment?limit=25&start=25"},
					}},
				}},
				Links: Links{Base: "http://" + URL + "/wiki"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := &url.URL{
				Scheme: "http",
				Host:   URL,
				Path:   tt.args.path,
			}
			got, err := api.SendContentRequest(url, tt.args.method, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendContentRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendContentRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
