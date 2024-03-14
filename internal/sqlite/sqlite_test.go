package sqlite

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gsiffert/fetch/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetaDataRepo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo, err := NewMetaDataRepo(ctx, "file:test.sqlite?cache=shared&mode=memory")
	require.NoError(t, err)

	defer func() {
		err := repo.Close()
		require.NoError(t, err)
	}()

	records := []domain.MetaData{
		{
			ID:          domain.PageID("https://wwww.google.com"),
			Site:        "www.google.com",
			LastFetched: time.Now().UTC().Truncate(time.Second),
			NumImages:   18,
			NumLinks:    8,
		},
		{
			ID:          domain.PageID("https://wwww.google.com/abount"),
			Site:        "www.google.com/about",
			LastFetched: time.Now().UTC().Truncate(time.Second),
			NumImages:   35,
			NumLinks:    23,
		},
	}

	for _, record := range records {
		t.Run(fmt.Sprintf("save metatdata %s", record.Site), func(t *testing.T) {
			err := repo.Save(ctx, record)
			require.NoError(t, err)
		})
	}

	t.Run("retrieve everything", func(t *testing.T) {
		fetchedRecords, err := repo.ByIDs(ctx, []domain.PageID{records[0].ID, records[1].ID})
		require.NoError(t, err)
		assert.Equal(t, records, fetchedRecords)
	})
}
