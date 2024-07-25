

all: build run


build:
	@go build -o ./bin/server ./cmd/main.go
	@ echo "Build a server"
run:
	@echo "Run server"
	@./bin/server

clean:
	@rm -rf ./bin/server
	@echo "Clean"

.PHONY: all build run clean