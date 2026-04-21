.PHONY: build install

build:
	GOOS=linux GOARCH=amd64 go build -o bin/copilot-cli-env-linux
	GOOS=windows GOARCH=amd64 go build -o bin/copilot-cli-env.exe
	GOOS=darwin GOARCH=arm64 go build -o bin/copilot-cli-env-mac

install:
	go install .