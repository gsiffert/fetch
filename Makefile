
.PHONY: vendor
vendor:
	go mod vendor

.PHONY: build-go
build-go: vendor
	CGO_ENABLED=1 go build ./cmd/fetch

.PHONY: build-docker
build-docker: vendor
	docker compose build fetch

.PHONY: unit-tests
unit-tests: vendor
	go test -count=1 ./...
