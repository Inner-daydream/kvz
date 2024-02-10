.PHONY: sqlgen

sqlgen:
	cd internal/sqlite/sql && sqlc generate
