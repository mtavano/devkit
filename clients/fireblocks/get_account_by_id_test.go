package fireblocks

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/mtavano/devkit/test"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_GetAccount(t *testing.T) {
	testCases := []struct {
		name               string
		baseURL            string
		apiKey             string
		AccountID          string
		mockHTTPClient     *mockHTTPClient
		expectedErrMsg     string
		expectedResponse   *VaultAccount
		generatePrivateKey func() []byte
	}{
		{
			name: "should return an error if privatekey is nil",
			mockHTTPClient: &mockHTTPClient{
				err: errors.New("expectedErr1"),
			},
			expectedErrMsg: "fireblocks: could not create getAccounts request: fireblocks: Client.makeRequest cl.signJWT error: fireblocks: Client.signJWT jwt.ParseRSAPrivateKeyFromPEM error: Invalid Key: Key must be a PEM encoded PKCS1 or PKCS8 key",
			generatePrivateKey: func() []byte {
				return nil
			},
		},
		{
			name: "should return an error if http client return an error",
			mockHTTPClient: &mockHTTPClient{
				err: errors.New("expectedErr1"),
			},
			expectedErrMsg: "fireblocks: could not create getAccounts request: fireblocks: Client.makeRequest cl.httpClient.Do error: expectedErr1",
			generatePrivateKey: func() []byte {
				privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					fmt.Printf("Cannot generate RSA key\n")
					os.Exit(1)
				}

				// dump private key to file
				privateKeyBytes := x509.MarshalPKCS1PrivateKey(privatekey)
				privateKeyBlock := &pem.Block{
					Type:  "RSA PRIVATE KEY",
					Bytes: privateKeyBytes,
				}
				return pem.EncodeToMemory(privateKeyBlock)
			},
		},

		{
			name:      "should return a successful response",
			AccountID: "1",
			mockHTTPClient: &mockHTTPClient{
				res: test.CreateMockResponse(`{
					"id": "1",
					"name": "Network Deposits",
					"hiddenOnUI": false,
					"autoFuel": false,
					"assets": [
						{
							"id": "BTC",
							"total": "0",
							"balance": "0",
							"lockedAmount": "0",
							"available": "0",
							"pending": "0",
							"frozen": "0",
							"staked": "0",
							"blockHeight": "-1"
						},
						{
							"id": "ETH",
							"total": "0.011465026322472291",
							"balance": "0.011465026322472291",
							"lockedAmount": "0",
							"available": "0.011465026322472291",
							"pending": "0",
							"frozen": "0",
							"staked": "0",
							"blockHeight": "13747745"
						},
						{
							"id": "USDC",
							"total": "27777.518154",
							"balance": "27777.518154",
							"lockedAmount": "0",
							"available": "27777.518154",
							"pending": "0",
							"frozen": "0",
							"staked": "0",
							"blockHeight": "13747745"
						},
						{
							"id": "TERRA_USD",
							"total": "0",
							"balance": "0",
							"lockedAmount": "0",
							"available": "0",
							"pending": "0",
							"frozen": "0",
							"staked": "0",
							"blockHeight": "-1"
						},
						{
							"id": "USDT_ERC20",
							"total": "0",
							"balance": "0",
							"lockedAmount": "0",
							"available": "0",
							"pending": "0",
							"frozen": "0",
							"staked": "0",
							"blockHeight": "13687126"
						}
					]
				}`, http.StatusOK),
			},
			expectedErrMsg: "",
			expectedResponse: &VaultAccount{
				ID:         "1",
				Name:       "Network Deposits",
				HiddenOnUI: false,
				AutoFuel:   false,
				Assets: []*VaultAccountAsset{
					{
						ID:           "BTC",
						Total:        "0",
						Balance:      "0",
						LockedAmount: "0",
						Available:    "0",
						Pending:      "0",
						Frozen:       "0",
						Staked:       "0",
						BlockHeight:  "-1",
					},
					{
						ID:           "ETH",
						Total:        "0.011465026322472291",
						Balance:      "0.011465026322472291",
						LockedAmount: "0",
						Available:    "0.011465026322472291",
						Pending:      "0",
						Frozen:       "0",
						Staked:       "0",
						BlockHeight:  "13747745",
					},
				},
			},
			generatePrivateKey: func() []byte {
				privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					fmt.Printf("Cannot generate RSA key\n")
					os.Exit(1)
				}

				// dump private key to file
				privateKeyBytes := x509.MarshalPKCS1PrivateKey(privatekey)
				privateKeyBlock := &pem.Block{
					Type:  "RSA PRIVATE KEY",
					Bytes: privateKeyBytes,
				}
				return pem.EncodeToMemory(privateKeyBlock)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fireblocksClient := NewClient(
				tc.baseURL,
				tc.apiKey,
				tc.mockHTTPClient,
			)

			fireblocksClient.LoadPrivateKey(tc.generatePrivateKey())

			res, err := fireblocksClient.GetAccount(tc.AccountID)
			if err != nil {
				require.Nil(t, res)
				require.Equal(t, tc.expectedErrMsg, err.Error())
				return
			}

			if res != nil {
				require.Nil(t, err)

				require.Equal(t, tc.expectedResponse.ID, res.ID)
				require.Equal(t, tc.expectedResponse.Name, res.Name)
				require.Equal(t, tc.expectedResponse.HiddenOnUI, res.HiddenOnUI)
				require.Equal(t, tc.expectedResponse.CustomerRefId, res.CustomerRefId)
				require.Equal(t, tc.expectedResponse.AutoFuel, res.AutoFuel)

				for id, asset := range tc.expectedResponse.Assets {
					require.Equal(t, asset.ID, res.Assets[id].ID)
					require.Equal(t, asset.Total, res.Assets[id].Total)
					require.Equal(t, asset.Balance, res.Assets[id].Total)
					require.Equal(t, asset.Available, res.Assets[id].Available)
					require.Equal(t, asset.Pending, res.Assets[id].Pending)
					require.Equal(t, asset.LockedAmount, res.Assets[id].LockedAmount)
					require.Equal(t, asset.Staked, res.Assets[id].Staked)
					require.Equal(t, asset.BlockHeight, res.Assets[id].BlockHeight)

				}
			}
		})
	}
}
