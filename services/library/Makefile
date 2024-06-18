.PHONY: up down start fix

up:
	migrate -path migrations/postgres/ -database "postgresql://postgres:root@localhost:5432/library?sslmode=disable" -verbose up

down:
	migrate -path migrations/postgres/ -database "postgresql://postgres:root@localhost:5432/library?sslmode=disable" -verbose down

start:
	go run cmd/library/main.go

fix:
	migrate -path migrations/postgres/ -database "postgresql://postgres:root@localhost:5432/library?sslmode=disable" force 000001


