lint:
	gocritic check ./...
	revive ./...
	golint ./...
	goconst ./...
	golangci-lint run
	go vet ./...
	staticcheck ./...
run-test:
	cd ./internal/test && go test -v .
