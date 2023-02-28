package fintoc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/mtavano/devkit/errors"
)

const (
	getRefreshIntentURL = "/refresh_intents?%s"
)

// CreateRefreshIntent ...
func (cl *Client) GetRefreshIntent() ([]*RefreshIntent, error) {
	urlValues := url.Values{}
	urlValues.Add("link_token", cl.linkToken)
	path := fmt.Sprintf(
		getRefreshIntentURL,
		urlValues.Encode(),
	)

	res, err := cl.makeRequest(http.MethodGet, path)
	if err != nil {
		return nil, errors.Wrap(err, "fintoc: Client.GetRefreshIntent cl.makeRequest error")
	}

	var refreshIntents []*RefreshIntent
	err = cl.scanBody(res, &refreshIntents)
	if err != nil {
		return nil, errors.Wrap(err, "fintoc: Client.GetRefreshIntent cl.scanBody error")
	}

	return refreshIntents, nil
}
