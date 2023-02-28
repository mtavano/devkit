package fireblocks

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mtavano/devkit/errors"
)

var (
	ErrFireblocksUnauthorized = errors.New("Unauthorized access")
)

type Client struct {
	baseURL    string
	apiKey     string
	privateKey []byte
	httpClient HTTPClient
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func NewClient(baseURL string, apiKey string, httpClient HTTPClient) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (cl *Client) LoadPrivateKey(privateKey []byte) {
	cl.privateKey = privateKey
}

func (cl *Client) LoadPrivateKeyFromBase64(privateString string) error {
	privateKey, err := base64.StdEncoding.DecodeString(privateString)
	if err != nil {
		return err
	}

	cl.privateKey = privateKey
	return nil
}

func (cl *Client) makeRequest(method string, path string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", cl.baseURL, path)
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "fireblocks: Client.makeRequest http.NewRequest error")
	}

	accessToken, err := cl.signJWT(path, body)
	if err != nil {
		return nil, errors.Wrap(err, "fireblocks: Client.makeRequest cl.signJWT error")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("X-API-Key", cl.apiKey)

	res, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fireblocks: Client.makeRequest cl.httpClient.Do error")
	}

	if err := cl.checkStatusCode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (cl *Client) checkStatusCode(res *http.Response) error {
	if res.StatusCode == http.StatusOK {
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "fireblocks: Client.checkStatusCode ioutil.ReadAll error")
	}
	defer res.Body.Close()

	errorMessage := fmt.Sprintf("fireblocks: unknown status code %d with body: %s", res.StatusCode, string(body))
	return errors.New(errorMessage)
}

func (cl *Client) scanBody(res *http.Response, scanner interface{}) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "fireblocks: Client.scanBody ioutil.ReadAll error")
	}
	defer res.Body.Close()

	return json.Unmarshal(body, scanner)
}

func (cl *Client) signJWT(path string, bodyJson []byte) (string, error) {
	type MyCustomClaims struct {
		Uri      string `json:"uri"`
		Nonce    string `json:"nonce"`
		BodyHash string `json:"bodyHash"`

		jwt.StandardClaims
	}

	var bodyHash string
	if bodyJson != nil {
		sha := sha256.New()
		sha.Write(bodyJson)
		bodyHash = hex.EncodeToString(sha.Sum(nil))
	}

	// Create the Claims
	claims := MyCustomClaims{
		Uri:      path,
		Nonce:    "test-nonce",
		BodyHash: bodyHash,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  int64(time.Now().Unix()),
			ExpiresAt: int64(time.Now().Unix()) + 55,
			Subject:   cl.apiKey,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(cl.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "fireblocks: Client.signJWT jwt.ParseRSAPrivateKeyFromPEM error")
	}

	return token.SignedString(signKey)
}

func (cl *Client) verify(tokenString string, publicKey []byte) (*jwt.Token, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})

	return token, errors.Wrap(err, "fireblocks: Client.verify jwt.Parse error")
}
