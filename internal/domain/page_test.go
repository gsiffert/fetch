package domain

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected Page
	}{
		{
			name:  "without path",
			input: "https://www.google.com",
			expected: Page{
				ID:           PageID("https://www.google.com"),
				Site:         "www.google.com",
				FileLocation: "www.google.com",
			},
		},
		{
			name:  "with path",
			input: "https://www.google.com/about",
			expected: Page{
				ID:           PageID("https://www.google.com/about"),
				Site:         "www.google.com/about",
				FileLocation: "www.google.com%2Fabout",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			u, err := url.Parse(test.input)
			require.NoError(t, err)

			page := NewPage(u)
			assert.Equal(t, test.expected, page)
		})
	}
}
