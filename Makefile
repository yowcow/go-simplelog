all:
	which dep || go get github.com/golang/dep/cmd/dep && go install github.com/golang/dep/cmd/dep
	dep ensure -v

test:
	go test

.PHONY: all test
