BINARY_NAME=pizzeria-app
CONTAINER_NAME=pizzeria-postgres
VOLUME_NAME=pizzeria-pgdata
SCHEMA_FILE=db/schema.sql
DB_URL=postgres://postgres:postgres@localhost:5432/pizzeria?sslmode=disable

.PHONY: all run-app clean db-up db-down db-clean schema-init sqlc-gen templ-gen build execute dev

# Regla principal para la entrega con manejador de salida y limpieza acotada al proyecto
run-app: db-clean db-up wait-db schema-init sqlc-gen templ-gen build
	@trap 'make cleanup-on-exit' EXIT INT TERM; \
	$(MAKE) execute

cleanup-on-exit:
	@echo "Ejecutando limpieza posterior del proyecto..."
	@docker stop $(CONTAINER_NAME) 2>/dev/null || true
	@docker rm -f $(CONTAINER_NAME) 2>/dev/null || true
	@docker volume rm -f $(VOLUME_NAME) 2>/dev/null || true

# --- 1) TAREAS PREVIAS Y GESTIÓN DE BD ---
db-clean:
	@echo "Limpiando contenedor y volumen del proyecto..."
	@docker stop $(CONTAINER_NAME) 2>/dev/null || true
	@docker rm -f $(CONTAINER_NAME) 2>/dev/null || true
	@docker volume rm -f $(VOLUME_NAME) 2>/dev/null || true

db-down: db-clean

db-up:
	@echo "Levantando contenedor de PostgreSQL con volumen propio del proyecto..."
	@docker volume create $(VOLUME_NAME) >/dev/null
	docker run --name $(CONTAINER_NAME) \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=pizzeria \
		-v $(VOLUME_NAME):/var/lib/postgresql/data \
		-p 5432:5432 \
		-d postgres:15-alpine

wait-db:
	@echo "Esperando a que la base de datos esté lista y aceptando conexiones..."
	@until docker exec $(CONTAINER_NAME) pg_isready -U postgres -d pizzeria >/dev/null 2>&1 && \
		docker exec $(CONTAINER_NAME) psql -U postgres -d pizzeria -c "SELECT 1;" >/dev/null 2>&1; do \
		sleep 1; \
	done
	@sleep 1
	@echo "¡Base de datos lista!"

schema-init:
	@echo "Aplicando esquema desde $(SCHEMA_FILE)..."
	@for i in 1 2 3 4 5; do \
		if docker exec -i $(CONTAINER_NAME) psql -U postgres -d pizzeria < $(SCHEMA_FILE); then \
			break; \
		fi; \
		echo "Reintentando aplicar esquema ($$i/5)..."; \
		sleep 2; \
	done

sqlc-gen:
	@echo "Ejecutando sqlc..."
	sqlc generate

templ-gen:
	@echo "Generando código de componentes templ..."
	templ generate

build:
	@echo "Compilando aplicación Go..."
	DATABASE_URL="$(DB_URL)" go build -o $(BINARY_NAME) .

# --- 2) EJECUTAR LA APLICACIÓN ---
execute:
	@echo "Iniciando servidor SSR de la Pizzería en http://localhost:8080..."
	@DATABASE_URL="$(DB_URL)" ./$(BINARY_NAME)

# --- 3) TAREAS POSTERIORES Y MODO DESARROLLO ---
clean: db-clean
	@echo "Limpiando binarios y archivos generados..."
	rm -f $(BINARY_NAME)
	find views -name "*_templ.go" -delete

# Opcional para desarrollo diario en tu máquina
dev: db-up wait-db schema-init
	@echo "Iniciando en modo desarrollo..."
	templ generate --watch & go run .