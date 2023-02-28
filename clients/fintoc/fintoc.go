package fintoc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mtavano/devkit/errors"
)

type Client struct {
	baseURL       string
	apiSecret     string
	accountID     string
	linkToken     string
	localCurrency string
	httpClient    HTTPClient

	accounts []*Account
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(baseURL, apiSecret, linkToken, localCurrency string, httpClient HTTPClient) *Client {

	return &Client{
		baseURL:       baseURL,
		apiSecret:     apiSecret,
		linkToken:     linkToken,
		localCurrency: localCurrency,
		httpClient:    httpClient,
	}
}

func (cl *Client) GetLocalCurrency() string {
	return cl.localCurrency
}

func (cl *Client) makeRequest(method, path string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", cl.baseURL, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "fintoc: Client.makeRequest http.NewRequest error")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", cl.apiSecret)

	res, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fintoc: Client.makeRequest cl.httpClient.Do error")
	}

	return res, nil
}

func (cl *Client) scanBody(res *http.Response, scanner interface{}) error {
	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated:
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, "fintoc: Client.scanBody ioutil.ReadAll error")
		}

		return json.Unmarshal(body, scanner)
	case http.StatusInternalServerError:
		return errors.New("internal server error")
	case http.StatusBadRequest:
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, "fintoc: Client.scanBody ioutil.ReadAll error")
		}

		return errors.Wrap(errors.New("bad request"), fmt.Sprintf("body: %s", string(body)))
	default:
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, "fintoc: Client.scanBody ioutil.ReadAll error")
		}

		return errors.New(fmt.Sprintf("unsupported status code %d with body: %s", res.StatusCode, string(body)))
	}
}
