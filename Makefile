.PHONY: all test build test_coverage

help:
	@echo '    build ...................... go get the dependencies'
	@echo '    run ........................ runs main.go'
	@echo '    test ....................... runs tests locally'
	@echo '    test_coverage .............. runs tests and generates coverage profile'


build:
	go get
	go get gopkg.in/mgo.v2
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	cd $$GOPATH/src/github.com/snickers/hls && make clean && make dep;
	go build

run: build
	DYLD_LIBRARY_PATH=$$GOPATH/src/github.com/snickers/hls/build ./snickers

test:
	@go vet ./...
	@ginkgo -r --slowSpecThreshold=20 --succinct .

test_coverage:
	@go get github.com/modocache/gover
	@ginkgo -r --slowSpecThreshold=20 --cover --succinct .
	@gover
	@mv gover.coverprofile coverage.txt

lint:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	gometalinter -disable=dupl ./... --deadline 300s
