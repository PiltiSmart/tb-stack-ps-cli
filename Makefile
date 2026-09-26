BINARY_NAME=pilti
BUILD_DIR=bin

.PHONY: all build install clean

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go mod tidy
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) main.go

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin/$(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installed successfully to /usr/local/bin/pilti"

clean:
	rm -rf $(BUILD_DIR)
