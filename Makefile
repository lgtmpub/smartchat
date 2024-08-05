.PHONY: all test benchmark dependencies lint vet fmt tidy osv-scan

all: osv-scan tidy test benchmark

test: dependencies lint vet fmt
	go test ./... -race -v -coverprofile cover.out

benchmark:
	go test -bench=. ./...

dependencies:
	go mod download

lint:
	@golangci-lint run

vet:
	go vet ./...

fmt:
	# @go install golang.org/x/tools/cmd/goimports@latest
	gofmt -l -s -w ./
	goimports -l -w ./

tidy:
	rm -f go.sum; go mod tidy -compat=1.22

osv-scan:
	# @go install github.com/google/osv-scanner/cmd/osv-scanner@v1
	@osv-scanner -r .