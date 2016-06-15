.PHONY: all test build

build:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega

test:
	ginkgo tests
