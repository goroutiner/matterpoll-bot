run:
	@echo "Запуск бота в демонстрационном режиме:"
	@docker compose up -d

stop:
	@echo "Остановка бота:"
	@docker compose down

fmt:
	@go fmt ./...

vet:
	@go vet ./...

unit-tests: vet
	@echo "Запуск unit-тестов для основной логики сервиса:"
	@go test -v ./internal/services/...

	@echo "Запуск unit-тестов для storage:"
	@go test -v internal/storage/validate_poll_test.go
	@go test -v ./internal/storage/memory/...

	@echo "Запуск unit-тестов для handlers:"
	@go test internal/handlers/middleware-unit_test.go

# integration-тесты запускаются только при запущенном Docker
integration-tests: unit-tests
	@echo "Запуск integration-тестов для storage:"
	@go test -v  internal/storage/database/database-integration_test.go

test-cover:
	@go test -cover ./...

clean:
	@go clean -testcache