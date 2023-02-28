package http

import (
	"context"
	nativehttp "net/http"

	"github.com/mtavano/devkit/errors"
	"golang.org/x/time/rate"
)

type BaseHTTPClient interface {
	Do(*nativehttp.Request) (*nativehttp.Response, error)
}

type ValidatorFunc func(*nativehttp.Response) (*nativehttp.Response, error)

type Client struct {
	httpClient BaseHTTPClient
	rl         *rate.Limiter

	validate ValidatorFunc
}

type Options struct {
	MaxRequest      int
	WindowInSeconds int
}

func NewClient(opts *Options, client BaseHTTPClient) *Client {
	rl := rate.NewLimiter(rate.Limit(opts.MaxRequest), opts.WindowInSeconds)

	return &Client{
		httpClient: client,
		rl:         rl,
	}
}

func (cl *Client) Do(req *nativehttp.Request) (res *nativehttp.Response, err error) {
	ctx := context.Background()

	// This is a blocking call
	err = cl.rl.Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "http: Client.Do cl.rl.Wait error")
	}
	res, err = cl.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "http: Client.Do endpoint[%s]", req.URL.EscapedPath())
	}

	if cl.validate != nil {
		res, err = cl.validate(res)
		if err != nil {
			// returned validation error
			return nil, err
		}
	}

	return res, nil
}

func (cl *Client) RegisterValidate(fn ValidatorFunc) {
	cl.validate = fn
}
