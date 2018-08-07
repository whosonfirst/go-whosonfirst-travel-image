CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-travel
	# cp *.go src/github.com/whosonfirst/go-whosonfirst-travel/
	# cp -r utils src/github.com/whosonfirst/go-whosonfirst-travel/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-image"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-travel"
	mv src/github.com/whosonfirst/go-whosonfirst-travel/vendor/github.com/whosonfirst/go-whosonfirst-geojson-v2 src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-travel/vendor/github.com/whosonfirst/go-whosonfirst-cli src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-travel/vendor/github.com/whosonfirst/go-whosonfirst-readwrite src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-image/vendor/github.com/whosonfirst/go-whosonfirst-flags src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/whosonfirst/go-whosonfirst-spr src/github.com/whosonfirst/
	rm -rf src/github.com/whosonfirst/go-whosonfirst-image/vendor/github.com/whosonfirst/go-whosonfirst-geojson-v2
	rm -rf src/github.com/whosonfirst/go-whosonfirst-image/vendor/github.com/whosonfirst/go-whosonfirst-spr
	rm -rf src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/whosonfirst/go-whosonfirst-flags

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	# go fmt *.go
	go fmt cmd/*.go

bin: 	self
	rm -rf bin/*
	@GOPATH=$(GOPATH) go build -o bin/wof-travel-id cmd/wof-travel-id.go
