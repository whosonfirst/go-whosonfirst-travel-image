cli:
	go build -mod vendor -o bin/wof-travel-id-image cmd/wof-travel-id-image/main.go
	go build -mod vendor -o bin/wof-belongs-to-image cmd/wof-belongs-to-image/main.go
	go build -mod vendor -o bin/wof-travel-filename cmd/wof-travel-filename/main.go
