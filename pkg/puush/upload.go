package puush

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

func (c *Client) Upload(file io.Reader, filename string) (string, error) {
	if !c.Account.Credentials.HasApiKey() {
		return "", PuushErrorInvalidCredentials
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("f", filename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.WriteField("k", *c.Account.Credentials.Key)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.FormatURL("/api/up"), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "puush")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	scanner, err := c.EvaluateResponse(resp)
	if err != nil {
		return "", err
	}

	responseData := strings.Split(scanner.Text(), ",")
	uploadUrl := responseData[1]
	updatedDiskUsage, err := strconv.ParseInt(responseData[2], 10, 64)
	if err != nil {
		return "", errors.New("response error: invalid disk usage provided")
	}

	c.Account.DiskUsage = updatedDiskUsage
	return uploadUrl, nil
}
