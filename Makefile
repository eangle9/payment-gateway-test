migrate-down:
	- migrate -database cockroachdb://root@localhost:26257/payment?sslmode=disable -path internal/constant/query/schemas -verbose down
migrate-up:
	- migrate -database cockroachdb://root@localhost:26257/payment?sslmode=disable -path internal/constant/query/schemas -verbose up
migrate-create:
	- migrate create -ext sql -dir internal/constant/query/schemas -tz "UTC" $(args)
swagger:
	-swag fmt && swag init -g initiator/initiator.go
tests:
	- go test ./...  -count=1
air:
	- go install github.com/cosmtrek/air@latest
sqlc:
	- sqlc generate -f ./config/sqlc.yaml
lint:
	- golangci-lint run ./...
startdocker:
	sudo docker start stockshow2f
dev:
	go run cmd/main.go 
mg-up:
	docker-compose -f ./docker-compose_loc.yml up -d
mg-down:
	docker-compose -f ./docker-compose_loc.yml down
