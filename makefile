.PHONY: sqlc

sqlgen:
	cd repositories/sqlite/sql && sqlc generate
