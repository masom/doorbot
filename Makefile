# http://zduck.com/2014/go-project-structure-and-dependencies/

.PHONY: activate build doc fmt lint run test vendor_clean vendor_get vendor_update vet

program = doorbot
workdir = $(shell pwd)
path = $(GOPATH)/src/github.com/masom/doorbot

default:  build

build: vet
	cd $(path); go build -v -o ./bin/$(program)

doc:
	godoc -http:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources 
fmt:
	cd $(path); go fmt ./doorbot/...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	cd $(path); golint ./doorbot

run:
	cd $(path); go run main.go

test:
	cd $(path); go test ./doorbot/... -tags tests

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	cd $(path); go vet ./doorbot/...
