APP = capcut-toolbox
BUILD_DIR = .
VERSION ?= dev

build:
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(APP).exe .

clean:
	rm -f $(BUILD_DIR)/$(APP).exe
