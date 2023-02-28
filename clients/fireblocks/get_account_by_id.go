package fireblocks

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	getAccountByIdURL = "/v1/vault/accounts/%s"
)

func (cl *Client) GetAccount(id string) (*VaultAccount, error) {
	resp, err := cl.makeRequest(http.MethodGet, fmt.Sprintf(getAccountByIdURL, id), nil)
	if err != nil {
		return nil, errors.Wrap(err, "fireblocks: could not create getAccounts request")
	}
	var vault VaultAccount
	if err := cl.scanBody(resp, &vault); err != nil {
		return nil, errors.Wrap(err, "fireblocks: cl.scanBody: could not scan body into struct")
	}

	return &vault, nil
}
