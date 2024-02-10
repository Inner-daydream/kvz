MAIN_PACKAGE_PATH := ./cmd/
BINARY_NAME :=kvz

.PHONY: sqlgen
sqlgen:
	cd internal/sqlite/sql && sqlc gener

build:
	go build -o ${BINARY_NAME} ${MAIN_PACKAGE_PATH}/main.go

test:
	go test ./...

clean:
	go clean
	rm ${BINARY_NAME}
	rm kv.db