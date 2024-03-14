package disk

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPageWriter(t *testing.T) {
	t.Parallel()

	temporyDir := os.TempDir()
	client := New(temporyDir)
	expectedContent := "Hello World"

	t.Run("write to file", func(t *testing.T) {
		writer, err := client.NewPageWriter(context.Background(), "www.google.com")
		require.NoError(t, err)
		_, err = fmt.Fprint(writer, expectedContent)
		require.NoError(t, err)
		err = writer.Close()
		require.NoError(t, err)
	})

	t.Run("read content", func(t *testing.T) {
		expectedFile := filepath.Join(temporyDir, "www.google.com.html")
		file, err := os.Open(expectedFile)
		require.NoError(t, err)
		content, err := io.ReadAll(file)
		require.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	})
}
