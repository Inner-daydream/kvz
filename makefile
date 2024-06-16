MAIN_PACKAGE_PATH := ./cmd
SUBDIRS := $(shell find $(MAIN_PACKAGE_PATH) -mindepth 1 -maxdepth 1 -type d)
BINS := $(patsubst $(MAIN_PACKAGE_PATH)/%,%,$(SUBDIRS))

.PHONY: sqlgen populate build test clean $(BINS)

sqlgen:
	cd internal/sqlite/sql && sqlc generate

dev: build populate

populate:
	./scripts/populate.sh

build: $(BINS)

$(BINS):
	go build -o $@ $(MAIN_PACKAGE_PATH)/$@/main.go

test:
	go test ./...

clean:
	go clean
	rm -f $(BINS); rm -f kv.db;