.PHONY: all test build test_coverage

build:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega

test:
	ginkgo tests

test_coverage:
	go get github.com/modocache/gover
	go get github.com/mattn/goveralls
	ginkgo -cover -coverpkg=./... tests
	gover
	goveralls -service drone.io -coverprofile=gover.coverprofile -repotoken $(COVERALLS_TOKEN)
