#!/usr/bin/env bash
set -e

CONTAINER_NAME="postgres-pizzeria-test"

echo "=== [1/5] Limpiando procesos y contenedores previos ==="
docker stop $CONTAINER_NAME 2>/dev/null || true
docker rm -f $CONTAINER_NAME 2>/dev/null || true

echo "=== [2/5] Compilando aplicación Go ==="
if command -v sqlc &> /dev/null; then
    sqlc generate
fi
go build -o pizzeria-api main.go

echo "=== [3/5] Levantando contenedor de PostgreSQL ==="
docker run -d \
  --name $CONTAINER_NAME \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=pizzeria \
  -p 5432:5432 \
  postgres:15-alpine

echo "=== [4/5] Esperando a que PostgreSQL esté listo... ==="
# Esperar a que el proceso pg_isready responda OK
until docker exec $CONTAINER_NAME pg_isready -h localhost -U postgres -d pizzeria &> /dev/null; do
    echo "Esperando conexión a PostgreSQL..."
    sleep 2
done

# Le damos 2 segundos extra para asegurar la inicialización completa del script de postgres
sleep 2
echo "Base de datos disponible."

# Cargar el esquema SQL usando el host y la base de datos explícitos
docker exec -i $CONTAINER_NAME psql -h localhost -U postgres -d pizzeria < db/schema.sql

echo "=== [5/5] Levantando la API REST ==="
DATABASE_URL="postgres://postgres:postgres@localhost:5432/pizzeria?sslmode=disable" ./pizzeria-api &
APP_PID=$!
sleep 3

# Función de limpieza al finalizar
cleanup() {
    echo "=== Tareas posteriores: Limpiando procesos y contenedor ==="
  if [ -n "$APP_PID" ] && kill -0 $APP_PID 2>/dev/null; then
        kill $APP_PID 2>/dev/null || true
        wait $APP_PID 2>/dev/null || true
    fi
    
    # Detener y eliminar contenedor y limpiar volúmenes anónimos asociados
    docker stop $CONTAINER_NAME 2>/dev/null || true
    docker rm -f $CONTAINER_NAME 2>/dev/null || true
    docker volume prune -f &>/dev/null || true
    rm -f pizzeria-api
    echo "=== Fin de la limpieza ==="
}
trap cleanup EXIT

echo "=== Ejecutando suite de pruebas HURL ==="
hurl --test requests.hurl

echo "=== ¡Todas las pruebas de la API se completaron con éxito! ==="