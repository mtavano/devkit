package fireblocks

import "net/http"

type mockHTTPClient struct {
	res *http.Response
	err error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.res, m.err
}
