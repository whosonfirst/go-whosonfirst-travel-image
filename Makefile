CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-travel-image
	# cp *.go src/github.com/whosonfirst/go-whosonfirst-travel-image/
	if test -d assets; then cp -r assets src/github.com/whosonfirst/go-whosonfirst-travel-image/; fi
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-image"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-travel"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-bindata"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-bindata-html-template"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/warning"
	@GOPATH=$(GOPATH) go get -u "golang.org/x/image/..."
	mv src/github.com/whosonfirst/go-whosonfirst-travel/vendor/github.com/whosonfirst/go-whosonfirst-geojson-v2 src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-travel/vendor/github.com/whosonfirst/go-whosonfirst-cli src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-travel/vendor/github.com/whosonfirst/go-whosonfirst-readwrite src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-image/vendor/github.com/whosonfirst/go-whosonfirst-flags src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/whosonfirst/go-whosonfirst-spr src/github.com/whosonfirst/
	rm -rf src/github.com/whosonfirst/go-whosonfirst-image/vendor/github.com/whosonfirst/go-whosonfirst-geojson-v2
	rm -rf src/github.com/whosonfirst/go-whosonfirst-image/vendor/github.com/whosonfirst/go-whosonfirst-spr
	rm -rf src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/whosonfirst/go-whosonfirst-flags
	rm -rf src/github.com/whosonfirst/go-bindata/testdata/

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	# go fmt *.go
	go fmt cmd/*.go

assets: self
	@GOPATH=$(GOPATH) go build -o bin/go-bindata ./vendor/github.com/whosonfirst/go-bindata/go-bindata/
	rm -rf templates/*/*~
	if test -d assets; then rm -rf assets; fi
	mkdir -p assets/html
	@GOPATH=$(GOPATH) bin/go-bindata -pkg html -o assets/html/html.go templates/html

bin: 	self
	rm -rf bin/*
	@GOPATH=$(GOPATH) go build -o bin/wof-travel-id cmd/wof-travel-id.go
