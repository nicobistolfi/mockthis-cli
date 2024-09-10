.PHONY: run build build-linux build-windows build-mac
run:
	go run cmd/mockthis/main.go
	
build:
	go build -o mockthis-cli cmd/mockthis/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/mockthis-linux-amd64 cmd/mockthis/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o ./bin/mockthis-windows-amd64.exe cmd/mockthis/main.go

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/mockthis-darwin-amd64 cmd/mockthis/main.go