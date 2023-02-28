package fintoc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/mtavano/devkit/errors"
)

const (
	createRefreshIntentURL = "/refresh_intents?%s"
)

// CreateRefreshIntent ...
func (cl *Client) CreateRefreshIntent() (*RefreshIntent, error) {
	urlValues := url.Values{}
	urlValues.Add("link_token", cl.linkToken)
	path := fmt.Sprintf(
		createRefreshIntentURL,
		urlValues.Encode(),
	)

	res, err := cl.makeRequest(http.MethodPost, path)
	if err != nil {
		return nil, errors.Wrap(err, "fintoc: Client.CreateRefreshIntent cl.makeRequest error")
	}

	var refreshIntent *RefreshIntent
	err = cl.scanBody(res, &refreshIntent)
	if err != nil {
		return nil, errors.Wrap(err, "fintoc: Client.CreateRefreshIntent cl.scanBody error")
	}

	return refreshIntent, nil
}
