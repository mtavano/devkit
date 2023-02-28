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

func Test_GetAccounts(t *testing.T) {
	testCases := []struct {
		name               string
		baseURL            string
		apiKey             string
		mockHTTPClient     *mockHTTPClient
		expectedErrMsg     string
		expectedResponse   []*VaultAccount
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
			name: "should return a successful response",
			mockHTTPClient: &mockHTTPClient{
				res: test.CreateMockResponse(`[
					{
						"id": "0",
						"name": "Default",
						"hiddenOnUI": false,
						"autoFuel": false,
						"assets": [
							{
								"id": "ETH",
								"total": "0",
								"balance": "0",
								"lockedAmount": "0",
								"available": "0",
								"pending": "0",
								"frozen": "0",
								"staked": "0"
							},
							{
								"id": "BTC",
								"total": "0",
								"balance": "0",
								"lockedAmount": "0",
								"available": "0",
								"pending": "0",
								"frozen": "0",
								"staked": "0"
							}
						]
					},
					{
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
					}
				]`, http.StatusOK),
			},
			expectedErrMsg: "",
			expectedResponse: []*VaultAccount{
				{
					ID:         "0",
					Name:       "Default",
					HiddenOnUI: false,
					AutoFuel:   false,
					Assets: []*VaultAccountAsset{
						{
							ID:           "ETH",
							Total:        "0",
							Balance:      "0",
							LockedAmount: "0",
							Available:    "0",
							Pending:      "0",
							Frozen:       "0",
							Staked:       "0",
						},
					},
				},
				{
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

			res, err := fireblocksClient.GetAccounts()
			if err != nil {
				require.Nil(t, res)
				require.Equal(t, tc.expectedErrMsg, err.Error())
				return
			}

			if res != nil {
				require.Nil(t, err)
				require.Equal(t, len(tc.expectedResponse), len(res))
				for id, account := range res {
					require.Equal(t, tc.expectedResponse[id].ID, account.ID)
					require.Equal(t, tc.expectedResponse[id].Name, account.Name)
					require.Equal(t, tc.expectedResponse[id].HiddenOnUI, account.HiddenOnUI)
					require.Equal(t, tc.expectedResponse[id].CustomerRefId, account.CustomerRefId)
					require.Equal(t, tc.expectedResponse[id].AutoFuel, account.AutoFuel)
					require.Equal(t, tc.expectedResponse[id].Assets[0].ID, account.Assets[0].ID)
					require.Equal(t, tc.expectedResponse[id].Assets[0].Total, account.Assets[0].Total)
					require.Equal(t, tc.expectedResponse[id].Assets[0].Balance, account.Assets[0].Total)
					require.Equal(t, tc.expectedResponse[id].Assets[0].Available, account.Assets[0].Available)
					require.Equal(t, tc.expectedResponse[id].Assets[0].Pending, account.Assets[0].Pending)
					require.Equal(t, tc.expectedResponse[id].Assets[0].LockedAmount, account.Assets[0].LockedAmount)
					require.Equal(t, tc.expectedResponse[id].Assets[0].Staked, account.Assets[0].Staked)
					require.Equal(t, tc.expectedResponse[id].Assets[0].BlockHeight, account.Assets[0].BlockHeight)

				}
			}
		})
	}
}
