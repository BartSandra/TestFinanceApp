run:
	@echo "Запуск приложения..."
	@sudo docker-compose up --build

up:
	@echo "Поднимаем контейнеры..."
	@sudo docker-compose up -d

down:
	@echo "Останавливаем контейнеры..."
	@sudo docker-compose down

migrate:
	@echo "Выполняем миграции..."
	@sudo docker-compose run --rm goose -dir db/migrations up

test:
	@echo "Запуск тестов..."
	@go test -v ./tests

clean:
	@echo "Очищаем контейнеры и ресурсы..."
	@sudo docker-compose down -v

build:
	@echo "Строим образы..."
	@sudo docker-compose build

logs:
	@echo "Просмотр логов..."
	@sudo docker-compose logs -f

restart:
	@echo "Перезапускаем контейнеры..."
	@sudo docker-compose restart

lint:
	golangci-lint run

style:
	@echo "Formatting Go code with gofmt..."
	@gofmt -w .
	@echo "Go code formatting complete."

