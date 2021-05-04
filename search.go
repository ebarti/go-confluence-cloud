package confluentcloud

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// getSearchEndpoint creates the correct api endpoint
func (a *api) getSearchEndpoint() (*url.URL, error) {
	return url.ParseRequestURI(a.endPoint.String() + "/search")
}

// GetSearchContentResults queries content using ContentQuery
func (a *api) GetSearchContentResults(query SearchContentQuery) (*SearchPageResults, error) {
	ep, err := a.getSearchEndpoint()
	if err != nil {
		return nil, err
	}
	ep.RawQuery = addSearchQueryParams(query).Encode()

	req, err := http.NewRequest(http.MethodGet, ep.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := a.Request(req)
	if err != nil {
		return nil, err
	}

	var content SearchPageResults
	err = json.Unmarshal(res, &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// addSearchQueryParams adds the defined query parameters
func addSearchQueryParams(query SearchContentQuery) *url.Values {
	data := url.Values{}
	if query.Cql != "" {
		data.Set("cql", query.Cql)
	}
	if len(query.CqlContext) != 0 {
		cqlContext, err := json.Marshal(query.CqlContext)
		if err != nil {
			return &data
		}
		data.Set("cqlcontext", string(cqlContext))
	}
	if query.Limit != 0 {
		data.Set("limit", strconv.Itoa(query.Limit))
	}
	if query.cursor != "" {
		data.Set("cursor", query.cursor)
	}
	if query.IncludeArchivedSpaces {
		data.Set("includeArchivedSpaces ", "true")
	}
	return &data
}
