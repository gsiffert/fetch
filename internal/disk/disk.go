// Package disk is part of the infrastructure layer and it implements the service.Disk interface using a local directory.
package disk

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
)

// Client of the disk package.
type Client struct {
	basePath string
}

// New instantiates a new Client.
func New(basePath string) *Client {
	return &Client{basePath: basePath}
}

// NewPageWriter creates a new file for the given name.
func (c *Client) NewPageWriter(_ context.Context, name string) (io.WriteCloser, error) {
	fileName := fmt.Sprintf("%s.html", name)
	filePath := path.Join(c.basePath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return file, nil
}
