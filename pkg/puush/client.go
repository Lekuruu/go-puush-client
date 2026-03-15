package puush

import (
	"bufio"
	"errors"
	"net/http"
	"strconv"
)

// Client is the main struct for interacting with the puush API.
type Client struct {
	Account *Account
	BaseURL string

	httpClient *http.Client
}

func NewClientFromCredentials(creds Credentials) (*Client, error) {
	account, err := NewAccountFromCredentials(creds)
	if err != nil {
		return nil, err
	}
	return &Client{
		Account:    account,
		BaseURL:    "https://puush.me",
		httpClient: http.DefaultClient,
	}, nil
}

func NewClientFromApiKey(apiKey string) *Client {
	creds := Credentials{Key: &apiKey}
	client, _ := NewClientFromCredentials(creds)
	return client
}

func NewClientFromLogin(email, password string) *Client {
	creds := Credentials{Email: &email, Password: &password}
	client, _ := NewClientFromCredentials(creds)
	return client
}

func (c *Client) SetCredentials(creds Credentials) error {
	if !creds.IsValid() {
		return errors.New("invalid credentials: either API key or login must be provided")
	}
	c.Account.Credentials = &creds
	return nil
}

func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

func (c *Client) FormatURL(path string) string {
	return c.BaseURL + path
}

// EvaluateResponse checks the response for errors and returns a scanner if the request was successful
func (c *Client) EvaluateResponse(response *http.Response) (*bufio.Scanner, PuushError) {
	// Check for server errors first
	if response.StatusCode >= http.StatusInternalServerError {
		return nil, PuushErrorRequestFailure
	}

	// Parse response body for error codes
	// The first line should usually contain a status code
	scanner := bufio.NewScanner(response.Body)

	if !scanner.Scan() {
		return nil, PuushErrorRequestFailure
	}

	responseLine := scanner.Text()
	statusCode, err := strconv.Atoi(responseLine)
	if err != nil {
		return nil, PuushErrorRequestFailure
	}

	if statusCode >= 0 {
		// We have a successful response
		return scanner, nil
	}

	switch statusCode {
	case -1:
		return nil, PuushErrorInvalidCredentials
	case -2:
		return nil, PuushErrorRequestFailure
	case -3:
		return nil, PuushErrorChecksumFailure
	case -4:
		return nil, PuushErrorInsufficientStorage
	default:
		return nil, PuushErrorUnknown
	}
}
