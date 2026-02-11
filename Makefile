APP_NAME=cli
DIST_DIR=dist

.PHONY: tidy fmt build run release clean

tidy:
	go mod tidy

fmt:
	gofmt -w main.go config.go cmd/docker.go cmd/ui.go

build:
	go build -o $(APP_NAME) .

run:
	go run .

release: clean
	mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe .

clean:
	rm -rf $(DIST_DIR) $(APP_NAME)