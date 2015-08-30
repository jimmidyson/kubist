.PHONY: all
all: build

.PHONY: deps
deps:
	go get github.com/tools/godep
	go get golang.org/x/tools/cmd/cover

.PHONY: build
build: deps *.go **/*.go
	godep go build -a github.com/fabric8io/kubist/...

.PHONY: test
test: deps build
	godep go test --test.short github.com/fabric8io/kubist/...

.PHONY: coverage
coverage: deps build
	scripts/coverage.sh

.PHONY: integration
integration: deps
	godep go test -v --timeout=30m github.com/fabric8io/kubist/...

.PHONY: pr
pr: test coverage integration
	scripts/check_gofmt.sh
