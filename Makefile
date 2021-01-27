assets: self
	@GOPATH=$(GOPATH) go build -o bin/go-bindata ./vendor/github.com/whosonfirst/go-bindata/go-bindata/
	rm -rf templates/*/*~
	if test -d assets; then rm -rf assets; fi
	mkdir -p assets/html
	@GOPATH=$(GOPATH) bin/go-bindata -pkg html -o assets/html/html.go templates/html

cli:
	go build -mod vendor -o bin/wof-travel-id-image cmd/wof-travel-id-image/main.go
	go build -mod vendor -o bin/wof-belongs-to-image cmd/wof-belongs-to-image/main.go
	go build -mod vendor -o bin/wof-travel-filename cmd/wof-travel-filename/main.go
