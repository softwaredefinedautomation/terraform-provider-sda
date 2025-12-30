package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Sign up - Create new tenant, return tenant token upon successful creation
func (c *Client) SignUp(auth AuthStruct) (*AuthResponse, error) {
	if auth.Username == "" || auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}
	rb, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ident/v1/tenant", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

// SignIn - Get a new token for user
func (c *Client) SignIn() (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("provide username and password")
	}
	rb, err := json.Marshal(c.Auth)
	if err != nil {
		return nil, err
	}

	log.Print("Sending Request to ", c.HostURL)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ident/v1/user/login", c.HostURL), bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}

	// Mandatory header for sending JSON data
	req.Header.Set("Content-Type", "application/json")

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/ident/v1/user/login", c.HostURL), strings.NewReader(string(rb)))
	// if err != nil {
	// 	return nil, err
	// }

	body, err := c.DoRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	log.Print(ar.IdToken)
	return &ar, nil
}

// SignIn - Get a new token for user
func (c *Client) GetUserTokenSignIn(auth AuthStruct) (*AuthResponse, error) {
	if auth.Username == "" || auth.Password == "" {
		return nil, fmt.Errorf("provide username and password")
	}
	rb, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ident/v1/user/login", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req, nil)
	if err != nil {
		return nil, errors.New("Unable to login")
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

// SignOut - Revoke the token for a user
func (c *Client) SignOut(authToken *string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ident/v1/user/logout", c.HostURL), strings.NewReader(string("")))
	if err != nil {
		return err
	}

	body, err := c.DoRequest(req, authToken)
	if err != nil {
		return err
	}

	if string(body) != "Signed out user" {
		return errors.New(string(body))
	}

	return nil
}
