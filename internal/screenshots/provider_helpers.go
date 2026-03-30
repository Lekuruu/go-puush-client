package screenshots

import (
	"errors"
	"os"
)

// temporaryReadCloser is an io.ReadCloser that deletes the underlying file when closed
type temporaryReadCloser struct {
	file *os.File
	path string
}

func (c *temporaryReadCloser) Read(p []byte) (int, error) {
	return c.file.Read(p)
}

func (c *temporaryReadCloser) Close() error {
	closeErr := c.file.Close()
	removeErr := os.Remove(c.path)

	if closeErr != nil {
		return closeErr
	}
	if removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
		return removeErr
	}
	return nil
}
