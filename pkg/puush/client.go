package puush

import (
	"bufio"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// Client is the main struct for interacting with the puush API.
type Client struct {
	Account *Account
	BaseURL string

	httpClient *http.Client
}

func NewClientFromCredentials(creds *Credentials) (*Client, error) {
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

func NewClientFromApiKey(email, apiKey string) *Client {
	creds := &Credentials{Identifier: stringOrNil(email), Key: stringOrNil(apiKey)}
	client, _ := NewClientFromCredentials(creds)
	return client
}

func NewClientFromLogin(email, password string) *Client {
	creds := &Credentials{Identifier: stringOrNil(email), Password: stringOrNil(password)}
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
	responseStatus := strings.SplitN(responseLine, ",", 1)[0]
	statusCode, err := strconv.Atoi(responseStatus)
	if err != nil {
		// Assuming we have a successful response here too
		return scanner, nil
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

// EvaluateHttpResponse returns a puush error based on the http status code of the response
func (c *Client) EvaluateHttpResponse(response *http.Response) PuushError {
	if response.StatusCode >= http.StatusInternalServerError {
		return PuushErrorRequestFailure
	}

	switch response.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return PuushErrorNotFound
	case http.StatusUnauthorized, http.StatusForbidden:
		return PuushErrorInvalidCredentials
	default:
		return PuushErrorUnknown
	}
}

func stringOrNil(s string) *string {
	if s != "" {
		return &s
	}
	return nil
}
