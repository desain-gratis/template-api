include .env

run-app:
	go run app/*.go -socket-address=localhost:9000 -config app/config/local.json
