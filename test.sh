#!/usr/bin/env bash
set -e

CONTAINER_NAME="postgres-tp2-test"
DB_USER="admin"
DB_PASS="1234"
DB_NAME="pizzeria"
DB_PORT="5432"

echo "=== 1. Tareas Previas ==="
echo "--> Deteniendo y eliminando contenedores/volúmenes antiguos si existen..."
docker rm -f $CONTAINER_NAME 2>/dev/null || true
docker volume prune -f 2>/dev/null || true

echo "--> Levantando contenedor de PostgreSQL..."
docker run --name $CONTAINER_NAME \
  -e POSTGRES_USER=$DB_USER \
  -e POSTGRES_PASSWORD=$DB_PASS \
  -e POSTGRES_DB=$DB_NAME \
  -p $DB_PORT:5432 \
  -d postgres:latest

echo "--> Esperando a que la base de datos esté lista..."
until docker exec $CONTAINER_NAME pg_isready -U $DB_USER > /dev/null 2>&1; do
  sleep 1
done

# Agregamos esta pausa para dar tiempo a que PostgreSQL complete su inicio interno
sleep 2

echo "--> Base de datos lista para recibir conexiones."

echo "--> Aplicando schema.sql..."
docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME < db/schema.sql

echo "--> Esperando a que la base de datos esté lista..."
until docker exec $CONTAINER_NAME pg_isready -U $DB_USER > /dev/null 2>&1; do
  sleep 1
done
echo "--> Base de datos lista para recibir conexiones."

echo "--> Aplicando schema.sql..."
docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME < db/schema.sql

echo "--> Generando código Go con sqlc..."
if command -v sqlc &> /dev/null; then
    sqlc generate
else
    docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate
fi

echo ""
echo "=== 2. Ejecución de Tests ==="
export DATABASE_URL="postgres://$DB_USER:$DB_PASS@localhost:$DB_PORT/$DB_NAME?sslmode=disable"
go test -v ./...

TEST_EXIT_CODE=$?

echo ""
echo "=== 3. Tareas Posteriores ==="
echo "--> Limpiando contenedores y volúmenes..."
docker stop $CONTAINER_NAME > /dev/null
docker rm -v $CONTAINER_NAME > /dev/null

echo "--> Proceso finalizado."
exit $TEST_EXIT_CODE