package clients

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Client -
type Client struct {
	HostURL      string
	HTTPClient   *http.Client
	AccessToken  string
	IdToken      string
	RefreshToken string
	ExpiresIn    int64
	Auth         AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires in"`
	TokenType    string `json:"token_type"`
	Token        string `json:"token"`
}

// NewClient -
func NewRestClient(host, username, password *string) (*Client, error) {
	log.Print("Creating new REST Client")
	c := Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		HostURL: *host,
	}
	
	// If username, password or host url are not provided, return empty client
	if username == nil || password == nil || host == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.IdToken = ar.IdToken
	c.AccessToken = ar.AccessToken
	c.RefreshToken = ar.RefreshToken
	c.ExpiresIn = ar.ExpiresIn

	return &c, nil
}

func (c *Client) DoRequest(req *http.Request, authToken *string) ([]byte, error) {
	token := c.IdToken
	log.Print(token)

	if authToken != nil {
		token = *authToken
	}

	req.Header.Set("Authorization", token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
