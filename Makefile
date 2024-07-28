

all: install-deps

# Проверка наличия docker и установка, если не установлен
install-docker:
	@if ! [ -x "$$(command -v docker)" ]; then \
		@echo "Docker is not installed. Installing Docker..."; \
		curl -fsSL https://get.docker.com -o get-docker.sh; \
		sh get-docker.sh; \
		rm get-docker.sh; \
	else \
		echo "Docker is already installed."; \
	fi

# Проверка наличия docker-compose и установка, если не установлен
install-docker-compose:
	@if ! [ -x "$$(command docker compose version)" ]; then \
		@echo "Docker Compose is not installed. Installing Docker Compose..."; \
		curl -L "https://github.com/docker/compose/releases/download/2.29.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose; \
		chmod +x /usr/local/bin/docker-compose; \
		ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose; \
	else \
		echo "Docker Compose is already installed."; \
	fi

# Сборка контейнеров и их поднятие
build:
	@echo "Docker containers build and up"
	docker compose up --build -d

# Установка зависимостей
install-deps: install-docker install-docker-compose

# Запуск контейнеров
up:
	@echo "Docker containers up"
	docker compose up -d

# Просмотр логов контейнеров
logs:
	@echo "Docker containers logs"
	docker compose logs -f --tail=100

# Остановка контейнеров
down:
	@echo "Docker containers down"
	docker compose down

# Чтение лона сообщений из Kafka
# если изменили TOPIC (worker.go), то нужно изменить и в Makefile --topic
kafka_message:
	@echo "Send message to Kafka"
	@docker exec -it kafka kafka-console-producer --broker-list kafka:9092 --topic message

# Очистка данных PostgreSQL и Kafka
clean_data:
	echo "Clean"
	@rm -rf ./internal/postgres/db/* ./data/kafka/*


.PHONY: all install-docker install-docker-compose install-deps build up logs down clean_data

# Команда для запуска внутри контейнера kafka для чтения полученных сообщений
# kafka-console-consumer --bootstrap-server kafka:9092 --topic message --from-beginning