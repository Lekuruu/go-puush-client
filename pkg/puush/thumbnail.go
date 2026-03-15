package puush

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) Thumbnail(id int) (io.ReadCloser, error) {
	if !c.Account.Credentials.HasApiKey() {
		return nil, PuushErrorInvalidCredentials
	}

	params := url.Values{}
	params.Add("k", *c.Account.Credentials.Key)
	params.Add("i", strconv.Itoa(id))

	resp, err := http.NewRequest("POST", c.FormatURL("/api/thumb"), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	resp.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp.Header.Set("User-Agent", "puush")

	response, err := c.httpClient.Do(resp)
	if err != nil {
		return nil, err
	}

	return response.Body, c.EvaluateHttpResponse(response)
}
