package confluentcloud

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// getContentIDEndpoint creates the correct api endpoint by given id
func (a *API) getContentIDEndpoint(id string) (*url.URL, error) {
	return url.ParseRequestURI(a.endPoint.String() + "/content/" + id)
}

// getContentEndpoint creates the correct api endpoint
func (a *API) getContentEndpoint() (*url.URL, error) {
	return url.ParseRequestURI(a.endPoint.String() + "/content/")
}

// getContentChildEndpoint creates the correct api endpoint by given id and type
func (a *API) getContentChildEndpoint(id string, t string) (*url.URL, error) {
	return url.ParseRequestURI(a.endPoint.String() + "/content/" + id + "/child/" + t)
}

// getContentGenericEndpoint creates the correct api endpoint by given id and type
func (a *API) getContentGenericEndpoint(id string, t string) (*url.URL, error) {
	return url.ParseRequestURI(a.endPoint.String() + "/content/" + id + "/" + t)
}

// GetContent queries content using ContentQuery
func (a *API) GetContent(query ContentQuery) (*Content, error) {
	ep, err := a.getContentEndpoint()
	if err != nil {
		return nil, err
	}
	ep.RawQuery = addContentQueryParams(query).Encode()

	req, err := http.NewRequest("GET", ep.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := a.Request(req)
	if err != nil {
		return nil, err
	}

	var content Content
	err = json.Unmarshal(res, &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// GetContentFromNext queries content using Links previously retrieved
func (a *API) GetContentFromNext(links Links) (interface{}, error) {
	if links.Base == "" || links.Next == "" {
		return nil, nil
	}
	nextUrl := links.Base + links.Next
	req, err := http.NewRequest("GET", nextUrl, nil)
	if err != nil {
		return nil, err
	}
	res, err := a.Request(req)
	if err != nil {
		return nil, err
	}

	var content Content
	err = json.Unmarshal(res, &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (a *API) GetAttachmentsFromResult(result Results, baseURL string) ([]Results, error) {
	next := result.Children.Attachment.Links.Next
	results := result.Children.Attachment.Results
	for {
		if next == "" {
			break
		}
		rawQuery := baseURL + result.Children.Attachment.Links.Next
		req, err := http.NewRequest("GET", rawQuery, nil)
		if err != nil {
			return nil, err
		}

		res, err := a.Request(req)
		if err != nil {
			return nil, err
		}

		var content Content
		err = json.Unmarshal(res, &content)
		if err != nil {
			return nil, err
		}
		results = append(results, content.Results...)
		next = content.Links.Next
	}

	return results, nil
}

// addContentQueryParams adds the defined query parameters
func addContentQueryParams(query ContentQuery) *url.Values {

	data := url.Values{}
	if len(query.Expand) != 0 {
		data.Set("expand", strings.Join(query.Expand, ","))
	}
	//get specific version
	if query.Version != 0 {
		data.Set("version", strconv.Itoa(query.Version))
	}
	if query.Limit != 0 {
		data.Set("limit", strconv.Itoa(query.Limit))
	}
	if query.OrderBy != "" {
		data.Set("orderby", query.OrderBy)
	}
	if query.PostingDay != "" {
		data.Set("postingDay", query.PostingDay)
	}
	if query.SpaceKey != "" {
		data.Set("spaceKey", query.SpaceKey)
	}
	if query.Start != 0 {
		data.Set("start", strconv.Itoa(query.Start))
	}
	if query.Status != "" {
		data.Set("status", query.Status)
	}
	if query.Title != "" {
		data.Set("title", query.Title)
	}
	if query.Trigger != "" {
		data.Set("trigger", query.Trigger)
	}
	if query.Type != "" {
		data.Set("type", query.Type)
	}
	return &data
}
