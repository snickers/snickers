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

run:
	go run main.go

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
	gometalinter ./... --deadline 300s
