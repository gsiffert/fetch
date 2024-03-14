package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/gsiffert/fetch/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type nopCloserWriter struct {
	io.Writer
}

func (nopCloserWriter) Close() error { return nil }

const htmlContent = `
<!DOCTYPE html>
<html>
	<head>
		<title>Google</title>
	</head>
	<body>
		<a href="https://www.google.com">Google</a>
		<img src="https://www.google.com/logo.png" alt="Google" />
		<img src="https://www.google.com/background.png" alt="Google Background" />
		<a href="https://www.google.com/about">About</a>
		<a href="https://www.google.com/contact">Contact</a>
		<a href="https://www.google.com/search">Search</a>
	</body>
</html>
`

func TestService_Fetch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		sites      []string
		setupMocks func(svcTest *serviceTest)
		assertErr  assert.ErrorAssertionFunc
	}{
		{
			name:      "no sites",
			sites:     []string{},
			assertErr: assert.NoError,
		},
		{
			name:  "fetcher failed",
			sites: []string{"https://www.google.com"},
			setupMocks: func(svcTest *serviceTest) {
				svcTest.fetcher.EXPECT().
					Fetch(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("fetcher failed"))
			},
			assertErr: assert.Error,
		},
		{
			name:  "disk failed",
			sites: []string{"https://www.google.com"},
			setupMocks: func(svcTest *serviceTest) {
				fetchedItem := &FetchedItem{
					Content: io.NopCloser(strings.NewReader("")),
				}

				svcTest.fetcher.EXPECT().
					Fetch(gomock.Any(), gomock.Any()).
					Return(fetchedItem, nil)
				svcTest.disk.EXPECT().
					NewPageWriter(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("disk failed"))
			},
			assertErr: assert.Error,
		},
		{
			name:  "save metadata failed",
			sites: []string{"https://www.google.com"},
			setupMocks: func(svcTest *serviceTest) {
				fetchedItem := &FetchedItem{
					Content: io.NopCloser(strings.NewReader("")),
				}
				writer := nopCloserWriter{io.Discard}

				svcTest.fetcher.EXPECT().
					Fetch(gomock.Any(), gomock.Any()).
					Return(fetchedItem, nil)
				svcTest.disk.EXPECT().
					NewPageWriter(gomock.Any(), gomock.Any()).
					Return(writer, nil)
				svcTest.metaDataRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(errors.New("save metadata failed"))
			},
			assertErr: assert.Error,
		},
		{
			name:  "success",
			sites: []string{"https://www.google.com"},
			setupMocks: func(svcTest *serviceTest) {
				fetchedItem := &FetchedItem{
					Page: domain.Page{
						ID:           domain.PageID("https://www.google.com"),
						Site:         "www.google.com",
						FileLocation: "www.google.com",
					},
					Content: io.NopCloser(strings.NewReader(htmlContent)),
				}
				writer := &bytes.Buffer{}
				writerCloser := nopCloserWriter{writer}

				before := time.Now().UTC()

				svcTest.fetcher.EXPECT().
					Fetch(gomock.Any(), string(fetchedItem.Page.ID)).
					Return(fetchedItem, nil)
				svcTest.disk.EXPECT().
					NewPageWriter(gomock.Any(), fetchedItem.Page.FileLocation).
					Return(writerCloser, nil)

				svcTest.metaDataRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Do(func(_ context.Context, m domain.MetaData) error {
						assert.Equal(t, 4, m.NumLinks)
						assert.Equal(t, 2, m.NumImages)
						assert.LessOrEqual(t, m.LastFetched, time.Now().UTC())
						assert.GreaterOrEqual(t, m.LastFetched, before)
						assert.Equal(t, fetchedItem.Page.Site, m.Site)
						assert.Equal(t, fetchedItem.Page.ID, m.ID)

						// We also verify that the writer received the content of the page.
						assert.Equal(t, htmlContent, writer.String())
						return nil
					})
			},
			assertErr: assert.NoError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			svcTest := newTestService(t)
			defer svcTest.Close()

			if test.setupMocks != nil {
				test.setupMocks(svcTest)
			}

			err := svcTest.svc.Fetch(ctx, test.sites...)
			test.assertErr(t, err)
		})
	}
}
