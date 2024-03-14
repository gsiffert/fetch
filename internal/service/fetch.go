package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gsiffert/fetch/internal/domain"
	"golang.org/x/net/html"
)

const (
	maxConcurrentFetch = 100
)

type FetchedItem struct {
	Page    domain.Page
	Content io.ReadCloser
}

func (f *FetchedItem) Close() error {
	return f.Content.Close()
}

// parseMetaData reads the html Content and returns the metadata.
func (s *Service) parseMetaData(_ context.Context, data io.Reader) (*domain.MetaData, error) {
	metaData := domain.MetaData{
		LastFetched: time.Now().UTC(),
	}

	// We use the html tokenizer instead of the parser to avoid parsing the whole document
	// as we only need to count the number of links and images.
	reader := html.NewTokenizer(data)
	for token := reader.Next(); token != html.ErrorToken; token = reader.Next() {
		if token != html.StartTagToken && token != html.SelfClosingTagToken {
			continue
		}

		tagName, _ := reader.TagName()
		tagNameStr := string(tagName)
		switch {
		case tagNameStr == "a":
			metaData.NumLinks++
		case tagNameStr == "img":
			metaData.NumImages++
		}
	}

	lastErr := reader.Err()
	if !errors.Is(lastErr, io.EOF) {
		return nil, fmt.Errorf("parse html: %w", lastErr)
	}

	return &metaData, nil
}

// fetchSite query the page, parse the metadata, saves the Content of the page in a file and save the metadata.
// The process stream the Content of the page to the file and through the metadata parser.
func (s *Service) fetchSite(ctx context.Context, site string) error {
	fetchedItem, err := s.fetcher.Fetch(ctx, site)
	if err != nil {
		return fmt.Errorf("query page: %w", err)
	}
	defer func() {
		if err := fetchedItem.Close(); err != nil {
			s.logger.Warn("Failed to close fetched item.", "site", site, "error", err)
		}
	}()

	writer, err := s.disk.NewPageWriter(ctx, fetchedItem.Page.FileLocation)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			s.logger.Warn("Failed to close writer.", "site", site, "error", err)
		}
	}()

	reader := io.TeeReader(fetchedItem.Content, writer)
	metaData, err := s.parseMetaData(ctx, reader)
	if err != nil {
		return fmt.Errorf("export metadata: %w", err)
	}

	metaData.ID = fetchedItem.Page.ID
	metaData.Site = fetchedItem.Page.Site
	if err := s.metaDataRepo.Save(ctx, *metaData); err != nil {
		return fmt.Errorf("save metadata: %w", err)
	}

	return nil
}

// fetchSitesInParallel fetches the sites in parallel and returns a channel of error.
// The channel will be closed when all the fetches are done.
func (s *Service) fetchSitesInParallel(ctx context.Context, sites []string) <-chan error {
	errChan := make(chan error)

	// We use a channel of empty structs to limit the number of concurrent fetches.
	pool := make(chan any, maxConcurrentFetch)

	go func() {
		defer close(errChan)

		// We run the fetch of each site in a goroutine and wait for all of them to finish.
		wg := sync.WaitGroup{}
		for _, site := range sites {
			pool <- nil
			wg.Add(1)

			go func(site string) {
				defer func() {
					<-pool
					wg.Done()
				}()

				err := s.fetchSite(ctx, site)
				if err != nil {
					select {
					case <-ctx.Done():
						return
					case errChan <- fmt.Errorf("fetch site %s: %w", site, err):
					}
				}
			}(site)
		}

		wg.Wait()
	}()

	return errChan
}

// Fetch downloads the sites, store their content in a file and save their related metadata.
// The sites are downloaded in parallel.
func (s *Service) Fetch(ctx context.Context, sites ...string) error {
	errChan := s.fetchSitesInParallel(ctx, sites)

	var errs error
	for err := range errChan {
		s.logger.Error("Failed to fetch site.", "error", err)
		errs = errors.Join(errs, err)
	}

	return errs
}
