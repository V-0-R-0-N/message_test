

all: install-deps build

# Проверка наличия docker и установка, если не установлен
install-docker:
	@if ! [ -x "$$(command -v docker)" ]; then \
		echo "Docker is not installed. Installing Docker..."; \
		curl -fsSL https://get.docker.com -o get-docker.sh; \
		sh get-docker.sh; \
		rm get-docker.sh; \
	else \
		echo -e "\nDocker is already installed.\n"; \
	fi

# Проверка наличия docker-compose и установка, если не установлен
install-docker-compose:
	@if ! [ -x "$$(command -v docker-compose)" ]; then \
		echo "Docker Compose is not installed. Installing Docker Compose..."; \
		curl -L "https://github.com/docker/compose/releases/download/2.29.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose; \
		chmod +x /usr/local/bin/docker-compose; \
		ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose; \
	else \
		echo -e "\nDocker Compose is already installed.\n"; \
	fi

# Сборка контейнеров и их поднятие
build:
	@echo -e "\nDocker containers build and up.\n"
	@docker compose up --build -d

# Установка зависимостей
install-deps: install-docker install-docker-compose

# Запуск контейнеров
up:
	@echo -e "\nDocker containers up.\n"
	@docker compose up -d

# Просмотр логов контейнеров
logs:
	@echo -e "\nDocker containers logs.\n"
	@docker compose logs -f --tail=100

# Остановка контейнеров
down:
	@echo -e "\nDocker containers down.\n"
	@docker compose down

# Чтение лона сообщений из Kafka
# если изменили TOPIC (worker.go), то нужно изменить и в Makefile --topic
kafka_message:
	@echo -e "\nLooking messages from Kafka\n"
	@docker compose exec kafka kafka-console-consumer --bootstrap-server kafka:9092 --topic message --from-beginning

# Очистка данных PostgreSQL и Kafka
clean_data:
	@echo -e "\nData is clean\n"
	@rm -rf ./data/postgres/* ./data/kafka/*


.PHONY: all install-docker install-docker-compose install-deps build up logs down kafka_message clean_data

# Команда для запуска внутри контейнера kafka для чтения полученных сообщений
# kafka-console-consumer --bootstrap-server kafka:9092 --topic message --from-beginning