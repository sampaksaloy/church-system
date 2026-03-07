.PHONY: run build deps tidy clean db-create db-drop help

APP_NAME=church-system
BINARY=./bin/server
DB_NAME=church_db

help:
	@echo "Church Event and Announcement System"
	@echo "────────────────────────────────────"
	@echo "make deps      - Download dependencies"
	@echo "make run       - Run the application"
	@echo "make build     - Build the binary"
	@echo "make tidy      - Tidy go modules"
	@echo "make db-create - Create the PostgreSQL database"
	@echo "make db-seed   - Seed sample data"
	@echo "make clean     - Remove build artifacts"

deps:
	go mod download

tidy:
	go mod tidy

run:
	go run ./cmd/server/main.go

build:
	mkdir -p bin
	go build -o $(BINARY) ./cmd/server/main.go
	@echo "Binary built at $(BINARY)"

db-create:
	createdb $(DB_NAME) || echo "Database may already exist"

db-drop:
	dropdb --if-exists $(DB_NAME)

db-seed:
	psql -d $(DB_NAME) -f db/seed.sql
	@echo "Sample data inserted"

clean:
	rm -rf bin/
