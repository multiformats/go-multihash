test: go_test

go_test: go_deps
	go test ./...

go_deps:
	go get code.google.com/p/go.crypto/sha3
	go get github.com/jbenet/go-base58
