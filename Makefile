.PHONY: build
build:
	@mkdir -p bin
	go build -trimpath \
		-ldflags "-s -w \
			-X github.com/cqroot/minop/pkg/version.version=dev \
			-X github.com/cqroot/minop/pkg/version.commit=$$(git rev-parse --short HEAD 2>/dev/null || echo none) \
			-X github.com/cqroot/minop/pkg/version.date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
		-o bin/minop .

.PHONY: test
test:
	go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: cover
cover: test
	go tool cover -html=coverage.out

.PHONY: fmt
fmt:
	gofumpt -w .

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: check
check:
	golangci-lint run
	@echo
	gofumpt -l .
