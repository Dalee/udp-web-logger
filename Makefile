.PHONY: install
install:
	go get github.com/UnnoTed/fileb0x
	go get github.com/golang/lint/golint
	go get github.com/gordonklaus/ineffassign
	go get github.com/client9/misspell/cmd/misspell
	fileb0x b0x.yaml

.PHONY: test
test:
	ineffassign ./
	gofmt -d -s -e main.go ./pkg/server/
	misspell -error README.md main.go ./pkg/server/*
