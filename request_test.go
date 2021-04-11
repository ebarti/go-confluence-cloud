package confluentcloud

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
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

	api, err := NewAPI(server.URL+"/wiki/rest/api", "userame", "token")
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

		req, err := http.NewRequest("GET", api.endPoint.String()+test.Endpoint, nil)
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
		case "/wiki/rest/api/content/":
			resp = Content{
				Results: []Results{Results{
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
