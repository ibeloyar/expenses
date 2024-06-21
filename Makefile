MAIN_FILE = cmd/expenses.go

run:
	go run $(MAIN_FILE)

install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest # swagger

swagger-gen:
	$(GOPATH)/bin/swag init -o ./api -g $(MAIN_FILE)