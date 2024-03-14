// Package fetcher is part of the infrastructure layer and it implements the service.Fetcher interface using a http.Client.
package fetcher

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/gsiffert/fetch/internal/domain"
	"github.com/gsiffert/fetch/internal/service"
)

const (
	maxRetries        = 5
	initialRetryDelay = 100 * time.Millisecond
	htmlContentType   = "text/html"
)

var retryErr = errors.New("retry")

// Client to Fetch webpages.
type Client struct {
	httpClient *http.Client
}

// New returns a new Client.
func New(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

// Fetch queries the page from the given site and returns a service.FetchedItem.
// It retries on network errors and server errors.
func (c *Client) Fetch(ctx context.Context, site string) (*service.FetchedItem, error) {
	r := retrier.New(retrier.ExponentialBackoff(maxRetries, initialRetryDelay), retrier.WhitelistClassifier{retryErr})

	var fetchedItem *service.FetchedItem
	err := r.Run(func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, site, nil)
		if err != nil {
			return fmt.Errorf("new request: %w", err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return errors.Join(fmt.Errorf("do request: %w", err), retryErr)
		}

		switch {
		case resp.StatusCode >= http.StatusInternalServerError:
			return errors.Join(fmt.Errorf("unexpected server error: %d", resp.StatusCode), retryErr)
		case resp.StatusCode != http.StatusOK:
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		case !strings.Contains(resp.Header.Get("Content-Type"), htmlContentType):
			return fmt.Errorf("unexpected content type: %s", resp.Header.Get("Content-Type"))
		default:
		}

		page := domain.NewPage(req.URL)
		fetchedItem = &service.FetchedItem{Page: page, Content: resp.Body}
		return nil
	})

	return fetchedItem, err
}
