package fireblocks

import (
	"net/http"

	"github.com/pkg/errors"
)

const (
	getAccountsURL = "/v1/vault/accounts"
)

func (cl *Client) GetAccounts() ([]*VaultAccount, error) {
	resp, err := cl.makeRequest(http.MethodGet, getAccountsURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "fireblocks: could not create getAccounts request")
	}
	var vault []*VaultAccount
	if err := cl.scanBody(resp, &vault); err != nil {
		return nil, errors.Wrap(err, "fireblocks: cl.scanBody: could not scan body into struct")
	}

	return vault, nil
}

type VaultAccount struct {
	// id	string	The ID of the Vault Account
	ID string `json:"id,omitempty"`
	// name	string	Name of the Vault Account
	Name string `json:"name,omitempty"`
	// hiddenOnUI	boolean	Specifies whether this vault account is visible in the web console or not
	HiddenOnUI bool `json:"hiddenOnUI,omitempty"`
	// customerRefId	string	[optional] The ID for AML providers to associate the owner of funds with transactions
	CustomerRefId string `json:"customerRefId,omitempty"`
	// autoFuel	boolean	Specifies whether this account's Ethereum address is auto fueled by the Gas Station or not
	AutoFuel bool `json:"autoFuel,omitempty"`
	// assets	array of VaultAccountAsset	List of assets under this Vault Account
	Assets []*VaultAccountAsset `json:"assets,omitempty"`
}

type VaultAccountAsset struct {
	// id	string	The ID of the asset
	ID string `json:"id,omitempty"`
	// total	string	The total wallet balance.
	Total string `json:"total,omitempty"`
	// balance	string	Deprecated - replaced by "total"
	Balance string `json:"balance,omitempty"`
	// available	string	Funds available for transfer. Equals the blockchain balance minus any locked amount
	Available string `json:"available,omitempty"`
	// pending	string	The cumulative balance of all transactions pending to be cleared
	Pending string `json:"pending,omitempty"`
	// staked	string	Staked funds, returned only for DOT
	Staked string `json:"staked,omitempty"`
	// frozen	string	Frozen by the AML policy in your workspace
	Frozen string `json:"frozen,omitempty"`
	// lockedAmount	string	Funds in outgoing transactions that are not yet published to the network
	LockedAmount string `json:"lockedAmount,omitempty"`
	// blockHeight	string	The height (number) of the block of the balance
	BlockHeight string `json:"blockHeight,omitempty"`
	// blockHash	string	The hash of the block of the balance
	BlockHash string `json:"blockHash,omitempty"`
}
