package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gsiffert/fetch/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_GetMetaDataForSites(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		sites      []string
		setupMocks func(svcTest *serviceTest)
		assertErr  assert.ErrorAssertionFunc
		expected   []domain.MetaData
	}{
		{
			name:      "no sites",
			sites:     []string{},
			assertErr: assert.NoError,
			setupMocks: func(svcTest *serviceTest) {
				svcTest.metaDataRepo.EXPECT().
					ByIDs(gomock.Any(), gomock.Any()).
					Return(nil, nil)
			},
		},
		{
			name:      "ByIDs failed",
			sites:     []string{},
			assertErr: assert.Error,
			setupMocks: func(svcTest *serviceTest) {
				svcTest.metaDataRepo.EXPECT().
					ByIDs(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("ByIDs failed"))
			},
		},
		{
			name:      "success",
			sites:     []string{"https://www.google.com", "https://www.google.com/about"},
			assertErr: assert.NoError,
			setupMocks: func(svcTest *serviceTest) {
				out := []domain.MetaData{
					{
						ID:          "https://www.google.com",
						Site:        "www.google.com",
						LastFetched: time.Date(2024, 3, 17, 14, 43, 0, 0, time.UTC),
						NumLinks:    12,
						NumImages:   2,
					},
					{
						ID:          "https://www.google.com/about",
						Site:        "www.google.com/about",
						LastFetched: time.Date(2024, 3, 17, 14, 43, 0, 0, time.UTC),
						NumLinks:    4,
						NumImages:   1,
					},
				}

				svcTest.metaDataRepo.EXPECT().
					ByIDs(gomock.Any(), []domain.PageID{out[0].ID, out[1].ID}).
					Return(out, nil)
			},
			expected: []domain.MetaData{
				{
					ID:          "https://www.google.com",
					Site:        "www.google.com",
					LastFetched: time.Date(2024, 3, 17, 14, 43, 0, 0, time.UTC),
					NumLinks:    12,
					NumImages:   2,
				},
				{
					ID:          "https://www.google.com/about",
					Site:        "www.google.com/about",
					LastFetched: time.Date(2024, 3, 17, 14, 43, 0, 0, time.UTC),
					NumLinks:    4,
					NumImages:   1,
				},
			},
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

			metadataItems, err := svcTest.svc.GetMetaDataForSites(ctx, test.sites...)
			test.assertErr(t, err)
			assert.Equal(t, test.expected, metadataItems)
		})
	}
}
