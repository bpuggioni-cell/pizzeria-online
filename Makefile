BINARY_NAME=pizzeria-app
CONTAINER_NAME=pizzeria-postgres
DB_URL=postgres://postgres:postgres@localhost:5432/pizzeria?sslmode=disable

.PHONY: all run-app clean db-up db-down db-clean sqlc-gen templ-gen build execute dev

# Regla principal para la entrega (ejecuta todo el ciclo completo de punta a punta)
run-app: db-clean db-up wait-db sqlc-gen templ-gen build execute db-clean

# --- 1) TAREAS PREVIAS Y GESTIÓN DE BD ---
db-clean:
	@echo "Limpiando contenedores y volúmenes anteriores de PostgreSQL..."
	@docker stop $(CONTAINER_NAME) 2>/dev/null || true
	@docker rm $(CONTAINER_NAME) 2>/dev/null || true
	@docker volume prune -f 2>/dev/null || true

db-up:
	@echo "Levantando contenedor de PostgreSQL con volumen limpio..."
	docker run --name $(CONTAINER_NAME) -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=pizzeria -p 5432:5432 -d postgres:15-alpine

wait-db:
	@echo "Esperando a que la base de datos esté lista y aceptando conexiones..."
	@until docker exec $(CONTAINER_NAME) pg_isready -U postgres; do \
		sleep 1; \
	done
	@echo "¡Base de datos lista!"

sqlc-gen:
	@echo "Ejecutando sqlc..."
	sqlc generate

templ-gen:
	@echo "Generando código de componentes templ..."
	templ generate

build:
	@echo "Compilando aplicación Go..."
	export DATABASE_URL="$(DB_URL)" && go build -o $(BINARY_NAME) .

# --- 2) EJECUTAR LA APLICACIÓN ---
execute:
	@echo "Iniciando servidor SSR de la Pizzería en http://localhost:8080..."
	export DATABASE_URL="$(DB_URL)" && ./$(BINARY_NAME)

# --- 3) TAREAS POSTERIORES Y MODO DESARROLLO ---
clean: db-clean
	@echo "Limpiando binarios y archivos generados..."
	rm -f $(BINARY_NAME)
	find views -name "*_templ.go" -delete

# Opcional para desarrollo diario en tu máquina
dev: db-up wait-db
	@echo "Iniciando en modo desarrollo..."
	templ generate --watch & go run .