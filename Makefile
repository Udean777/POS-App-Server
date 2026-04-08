include .env
export

.PHONY: db-reset dev tidy help

help:
	@echo "Tersedia perintah:"
	@echo "  make db-reset  - Menghapus dan membuat ulang database pos_db"
	@echo "  make dev       - Menjalankan server dengan live reload (air)"
	@echo "  make tidy      - Menjalankan go mod tidy"

db-reset:
	@echo "Sedang meriset database..."
	@psql -h localhost -U postgres -d postgres -c "DROP DATABASE IF EXISTS pos_db WITH (FORCE);"
	@psql -h localhost -U postgres -d postgres -c "CREATE DATABASE pos_db;"
	@echo "Database 'pos_db' berhasil diriset. Silakan jalankan server untuk migrasi ulang."

dev:
	@air

tidy:
	@go mod tidy
