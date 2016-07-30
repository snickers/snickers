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
	go vet ./...
	ginkgo -r -cover -keepGoing .

test_coverage:
	@go get github.com/modocache/gover
	@go get github.com/mattn/goveralls
	@ginkgo -cover -coverpkg=./... -r
	@gover
	@mv gover.coverprofile coverage.txt
