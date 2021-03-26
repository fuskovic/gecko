.PHONY: build lint zip

lint:
	goimports -w *.go
	@echo "go files linted"

build:
	@echo "removing old binary if one exists..."
	-rm gecko
	@echo "building new binary..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gecko

zip: build
	@echo "removing old zip file if one exists..."
	-rm gecko.zip
	@echo "creating new zip file..."
	zip gecko.zip gecko
	@echo "gecko.zip created"
	@rm gecko