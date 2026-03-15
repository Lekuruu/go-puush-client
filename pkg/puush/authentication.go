package puush

import (
	"errors"
	"net/http"
	"strings"
)

// Authenticate performs the authentication process using the client's credentials.
// The updated account data can then be accessed through the client's `Account` field.
// If authentication fails, an error is returned.
func (c *Client) Authenticate() error {
	if !c.Account.Credentials.IsValid() {
		return errors.New("invalid credentials: either API key or login must be provided")
	}
	params := c.Account.Credentials.toFormData()
	params.Add("z", "poop")
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", c.FormatURL("/api/auth"), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner, err := c.EvaluateResponse(resp)
	if err != nil {
		return err
	}

	accountIdentifer := c.Account.Credentials.Identifier
	account, err := NewAccountFromResponse(scanner)
	if err != nil {
		return errors.New("response error: " + err.Error())
	}

	c.Account = account
	c.Account.Credentials.Identifier = accountIdentifer
	return nil
}
