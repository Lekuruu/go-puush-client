package screenshots

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// temporaryFileReader is a `ReadSeekCloser` that deletes the underlying file when closed
type temporaryFileReader struct {
	file *os.File
	path string
}

func (c *temporaryFileReader) Read(p []byte) (int, error) {
	return c.file.Read(p)
}

func (c *temporaryFileReader) Seek(offset int64, whence int) (int64, error) {
	return c.file.Seek(offset, whence)
}

func (c *temporaryFileReader) Stat() (os.FileInfo, error) {
	return c.file.Stat()
}

func (c *temporaryFileReader) Close() error {
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

// temporaryDirectoryReader removes output directory when closed
type temporaryDirectoryReader struct {
	file *os.File
	dir  string
}

func (r *temporaryDirectoryReader) Read(p []byte) (int, error) {
	return r.file.Read(p)
}

func (r *temporaryDirectoryReader) Seek(offset int64, whence int) (int64, error) {
	return r.file.Seek(offset, whence)
}

func (r *temporaryDirectoryReader) Close() error {
	closeErr := r.file.Close()
	removeErr := os.RemoveAll(r.dir)
	if closeErr != nil {
		return closeErr
	}
	return removeErr
}

// idk bro i just needed a bytes reader that can close
type memoryReadCloser struct {
	*bytes.Reader
}

func (b *memoryReadCloser) Close() error {
	return nil
}

func findScreenshotInDirectory(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read screenshot directory: %w", err)
	}

	for _, entry := range entries {
		if entry.Type().IsRegular() {
			return filepath.Join(dir, entry.Name()), nil
		}
	}
	return "", errors.New("screenshot was cancelled or is empty")
}
