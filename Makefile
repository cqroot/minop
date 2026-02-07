.PHONY: test
test:
	mkdir -p '$(CURDIR)/pkg/remote/testdata'
	touch '$(CURDIR)/pkg/remote/testdata/test.txt'
	mkdir -p '$(CURDIR)/pkg/remote/testdata/test.dir'
	go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: cover
cover: test
	go tool cover -html=coverage.out

.PHONY: check
check:
	golangci-lint run
	@echo
	gofumpt -l .

