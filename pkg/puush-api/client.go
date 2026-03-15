package puush

import "errors"

type Client struct {
	Credentials *Credentials
	Account     *Account
	BaseURL     string
}

func NewClientFromCredentials(creds Credentials) (*Client, error) {
	if !creds.IsValid() {
		return nil, errors.New("invalid credentials provided")
	}
	return &Client{Credentials: &creds, BaseURL: "https://puush.me"}, nil
}

func NewClientFromApiKey(apiKey string) *Client {
	creds := Credentials{Key: &apiKey}
	client, _ := NewClientFromCredentials(creds)
	return client
}

func NewClientFromLogin(username, password string) *Client {
	creds := Credentials{Username: &username, Password: &password}
	client, _ := NewClientFromCredentials(creds)
	return client
}

func (c *Client) SetCredentials(creds Credentials) error {
	if !creds.IsValid() {
		return errors.New("invalid credentials provided")
	}
	c.Credentials = &creds
	return nil
}

func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}
