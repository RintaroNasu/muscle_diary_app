#!/bin/bash
set -euo pipefail

CONTAINER_NAME="${CONTAINER_NAME:-db}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-muscle_diary}"
EXERCISES=("ベンチプレス" "スクワット" "デッドリフト")

echo "Starting containers..."
docker compose up -d

running_containers () {
    local status
    status=$(docker compose ps -q | xargs -r docker inspect -f '{{.State.Status}}' || true)
    [[ -n "$status" && $(echo "$status" | grep -cv running) -eq 0 ]]
}

until running_containers; do sleep 1; done

echo "Waiting for PostgreSQL to be ready..."
until docker exec "$CONTAINER_NAME" pg_isready -U "$DB_USER" >/dev/null 2>&1; do
    sleep 2
done
echo "✅ PostgreSQL is ready."

echo "Seeding exercises..."
for exercise in "${EXERCISES[@]}"; do
docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" <<EOF
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM exercises WHERE name = '$exercise') THEN
        INSERT INTO exercises (name, created_at, updated_at)
        VALUES ('$exercise', NOW(), NOW());
        RAISE NOTICE 'Inserted: $exercise';
    ELSE
        RAISE NOTICE 'Skipped (already exists): $exercise';
    END IF;
END
\$\$;
EOF
done

echo "✅ Seed completed successfully!"
