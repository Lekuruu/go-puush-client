package puush

import (
	"errors"
	"net/http"
)

// Authenticate performs the authentication process using the client's credentials.
// The updated account data can then be accessed through the client's `Account` field.
// If authentication fails, an error is returned.
func (c *Client) Authenticate() error {
	if !c.Account.Credentials.IsValid() {
		return errors.New("invalid credentials: either API key or login must be provided")
	}
	params := c.Account.Credentials.toFormData()

	req, err := http.NewRequest("POST", c.FormatURL("/api/auth"), nil)
	if err != nil {
		return err
	}
	req.PostForm = params
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

	account, err := NewAccountFromResponse(scanner)
	if err != nil {
		return errors.New("response error: " + err.Error())
	}
	c.Account = account
	return nil
}
