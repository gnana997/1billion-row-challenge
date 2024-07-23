# Makefile for terminal_app

# Name of the binary
BINARY_NAME=1brc

# Go build command
BUILD_CMD=go build -o out/$(BINARY_NAME)

# Build the application
.PHONY: build
build:
	$(BUILD_CMD)
	@echo "Build completed."

# Run the application with create command
.PHONY: run-create
run-create:
	./out/$(BINARY_NAME) create --file=data/test.txt --rows=1000000000
	@echo "Create command executed."

# Run the application with read command
.PHONY: run-read
run-read:
	./out/$(BINARY_NAME) simple-process --file=data/test.txt
	@echo "Read command executed."

# Run the application with mmap command
.PHONY: run-mmap
run-mmap:
	./out/$(BINARY_NAME) use-basic-mmap --file=data/test.txt
	@echo "Mmap command executed."

# Run the application with parallel mmap command
.PHONY: run-parallel-mmap
run-parallel-mmap:
	./out/$(BINARY_NAME) use-parallel-mmap --file=data/test.txt
	@echo "Parallel mmap command executed."
# Clean the binary
.PHONY: clean
clean:
	rm -f out/$(BINARY_NAME)
	@echo "Clean completed."
