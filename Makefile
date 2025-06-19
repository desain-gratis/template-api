include .env

run:
	go run app/*.go -socket-address=localhost:9000 -config app/config/local.json

lint: 
	golangci-lint run app
