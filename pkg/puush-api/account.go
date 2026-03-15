package puush

import (
	"bufio"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Credentials represents the authentication credentials for a puush account.
// It can either contain an API key or a combination of email and password for login.
type Credentials struct {
	Email    *string
	Password *string
	Key      *string
}

func (c *Credentials) HasApiKey() bool {
	return c.Key != nil
}

func (c *Credentials) HasLoginCredentials() bool {
	return c.Email != nil && c.Password != nil
}

func (c *Credentials) IsValid() bool {
	return c.HasApiKey() || c.HasLoginCredentials()
}

func (c *Credentials) toFormData() url.Values {
	params := url.Values{}

	if c.HasApiKey() {
		params.Add("k", *c.Key)
	}
	if c.HasLoginCredentials() {
		params.Add("e", *c.Email)
		params.Add("p", *c.Password)
	}
	return params
}

// Account represents a puush account with its credentials, type, disk usage, and pro subscription end date.
// It can be obtained by calling `Authenticate()` on the Client struct, assuming the credentials are provided.
type Account struct {
	Credentials     *Credentials
	Type            AccountType
	DiskUsage       int64
	SubscriptionEnd *time.Time
}

func NewAccountFromCredentials(creds Credentials) (*Account, error) {
	if !creds.IsValid() {
		return nil, errors.New("invalid credentials: either API key or login must be provided")
	}
	return &Account{
		Credentials:     &creds,
		Type:            AccountTypeRegular,
		DiskUsage:       0,
		SubscriptionEnd: nil,
	}, nil
}

func NewAccountFromResponse(scanner *bufio.Scanner) (*Account, error) {
	if !scanner.Scan() {
		return nil, errors.New("failed to read account type from response")
	}
	authenticationResponse := strings.Split(scanner.Text(), ",")

	accountTypeInt, err := strconv.Atoi(authenticationResponse[0])
	if err != nil {
		return nil, errors.New("invalid account type in response")
	}
	accountType := AccountType(accountTypeInt)

	apiKey := authenticationResponse[1]
	subscriptionEndStr := authenticationResponse[2]

	var subscriptionEnd *time.Time

	if subscriptionEndStr != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", subscriptionEndStr)
		if err != nil {
			return nil, errors.New("invalid subscription end date format in response")
		}
		subscriptionEnd = &parsedTime
	}

	diskUsage, err := strconv.ParseInt(authenticationResponse[3], 10, 64)
	if err != nil {
		return nil, errors.New("invalid disk usage value in response")
	}

	return &Account{
		Type:            accountType,
		DiskUsage:       diskUsage,
		SubscriptionEnd: subscriptionEnd,
		Credentials:     &Credentials{Key: &apiKey},
	}, nil
}
