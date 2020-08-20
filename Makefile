VERSION_FILE=VERSION
VERSION=`cat $(VERSION_FILE)`

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build:
	rm -rf ./mysql-to-strict
	go build

.PHONY: build-docker-image
build-docker-image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o mysql-to-strict .
	docker build -t tomkeur/mysql-to-strict:$(VERSION) -f Dockerfile .
	rm -rf ./mysql-to-strict
