package test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func CreateMockResponse(body string, statusCode int) *http.Response {
	var bb io.ReadCloser
	if body != "" {
		bb = ioutil.NopCloser(bytes.NewBufferString(body))
	} else {
		bb = ioutil.NopCloser(bytes.NewBufferString("{}"))
	}

	return &http.Response{
		Body:       bb,
		StatusCode: statusCode,
	}
}
