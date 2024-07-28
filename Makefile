

all: build run


build:
	@go build -o ./bin/server ./cmd/main.go
	@ echo "Build a server"
run:
	@echo "Run server"
	@./bin/server

up:
	@echo "Docker containers up"
	docker-compose up -d

clean:
	@rm -rf ./bin/server
	@echo "Clean"

.PHONY: all build run clean

# Команда для запуска внутри контейнера kafka для чтения полученных сообщений
# kafka-console-consumer --bootstrap-server kafka:9092 --topic message --from-beginning