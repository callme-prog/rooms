.PHONY: up down test-auth test-booking swag-auth swag-booking tidy

## Запуск всего стека
up:
	docker compose up --build -d

## Остановка
down:
	docker compose down -v

## Тесты auth-service
test-auth:
	cd auth-service && go test ./... -v -count=1

## Тесты booking-service
test-booking:
	cd booking-service && go test ./... -v -count=1

## Все тесты с покрытием
test-cover:
	cd auth-service    && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out
	cd booking-service && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

## Генерация Swagger — auth-service
swag-auth:
	cd auth-service && swag init -g cmd/server/main.go -o docs

## Генерация Swagger — booking-service
swag-booking:
	cd booking-service && swag init -g cmd/server/main.go -o docs

## go mod tidy для всех модулей
tidy:
	cd common          && go mod tidy
	cd auth-service    && go mod tidy
	cd booking-service && go mod tidy
