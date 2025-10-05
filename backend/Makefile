# ===== Variables (override via env) =====
CONTAINER_NAME ?= db
DB_USER        ?= postgres
DB_NAME        ?= muscle_diary
# =======================================

.PHONY: seed up down ps logs migrate

up:
	@docker compose up -d
down:
	@docker compose down
ps:
	@docker compose ps
logs:
	@docker compose logs -f
seed:
	@echo "🌱 Seeding exercises (container=$(CONTAINER_NAME))..."
	@CONTAINER_NAME=$(CONTAINER_NAME) DB_USER=$(DB_USER) DB_NAME=$(DB_NAME) bash scripts/seed_exercises.sh
