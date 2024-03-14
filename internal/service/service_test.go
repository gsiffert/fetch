//go:generate mockgen -package service -source service.go -destination ./mock_test.go -typed
package service

import (
	"log/slog"
	"testing"

	"go.uber.org/mock/gomock"
)

type serviceTest struct {
	svc  *Service
	ctrl *gomock.Controller

	fetcher      *MockFetcher
	disk         *MockDisk
	metaDataRepo *MockMetaDataRepository
}

func (s *serviceTest) Close() {
	s.ctrl.Finish()
}

func newTestService(t *testing.T) *serviceTest {
	t.Helper()

	ctrl := gomock.NewController(t)
	svcTest := &serviceTest{
		ctrl:         ctrl,
		fetcher:      NewMockFetcher(ctrl),
		disk:         NewMockDisk(ctrl),
		metaDataRepo: NewMockMetaDataRepository(ctrl),
	}
	slog.SetLogLoggerLevel(slog.Level(10)) // Disable the logs.
	svcTest.svc = New(svcTest.fetcher, svcTest.disk, slog.Default(), svcTest.metaDataRepo)
	return svcTest
}
