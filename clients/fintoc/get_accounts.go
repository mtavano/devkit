package fintoc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/mtavano/devkit/errors"
)

const (
	getAccountsURL = "/accounts/?%s"
)

var (
	ErrAccountNotFound = errors.New("Account not found")
)

// GetAccountMovements will fetch the
func (cl *Client) GetAccounts() ([]*Account, error) {
	urlValues := url.Values{}
	urlValues.Add("link_token", cl.linkToken)

	path := fmt.Sprintf(
		getAccountsURL,
		urlValues.Encode(),
	)

	res, err := cl.makeRequest(http.MethodGet, path)
	if err != nil {
		return nil, errors.Wrap(err, "fintoc: Client.GetAccounts cl.makeRequest error")
	}

	accounts := make([]*Account, 0)
	err = cl.scanBody(res, &accounts)
	if err != nil {
		return nil, errors.Wrap(err, "fintoc: Client.GetAccounts cl.scanBody error")
	}

	cl.accounts = accounts

	return accounts, nil
}

func (cl *Client) GetAccountByType(accountType string) (*Account, error) {
	for _, acc := range cl.accounts {
		if accountType == acc.AccountType {
			return acc, nil
		}
	}

	return nil, ErrAccountNotFound
}

func (cl *Client) SetAccountIDOfType(accountType string) error {
	for _, acc := range cl.accounts {
		if accountType == acc.AccountType {
			cl.accountID = acc.ID
			return nil
		}
	}

	return ErrAccountNotFound
}
