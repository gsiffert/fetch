package fetcher

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gsiffert/fetch/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const htmlContent = `
<!DOCTYPE html>
<html>
	<head>
		<title>Test</title>
	</head>
	<body>
		<h1>Test</h1>
	</body>
</html>
`

type retryInternalErrorServer struct {
	count int
}

func (s *retryInternalErrorServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.count++
	if s.count < maxRetries-1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", htmlContentType)
	_, _ = w.Write([]byte(htmlContent))
}

func TestClient_Fetch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		server     http.Handler
		assertErr  assert.ErrorAssertionFunc
		expectResp bool
	}{
		{
			name: "invalid content type",
			server: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
			}),
			assertErr: assert.Error,
		},
		{
			name: "failed on 4xx",
			server: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}),
			assertErr: assert.Error,
		},
		{
			name: "failed on 5xx",
			server: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			assertErr: assert.Error,
		},
		{
			name: "success",
			server: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", htmlContentType)
				_, _ = w.Write([]byte(htmlContent))
			}),
			assertErr:  assert.NoError,
			expectResp: true,
		},
		{
			name:       "success after few retries",
			server:     &retryInternalErrorServer{},
			assertErr:  assert.NoError,
			expectResp: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			srv := httptest.NewServer(test.server)
			defer srv.Close()

			client := New(http.DefaultClient)
			item, err := client.Fetch(ctx, srv.URL)
			test.assertErr(t, err)
			if test.expectResp {
				withoutProtocol := strings.TrimPrefix(srv.URL, "http://")
				assert.Equal(t, domain.PageID(srv.URL), item.Page.ID)
				assert.Equal(t, withoutProtocol, item.Page.Site)
				assert.Equal(t, withoutProtocol, item.Page.FileLocation)
				b, err := io.ReadAll(item.Content)
				require.NoError(t, err)
				assert.Equal(t, htmlContent, string(b))
			}
		})
	}
}
