# Pizzería API REST - TP3

API REST desarrollada en Go (`net/http`) conectada a PostgreSQL mediante `sqlc` para la gestión de productos (Pizzas) y órdenes de compra (Pedidos).

## Requisitos Previos

* **Go** (versión 1.22 o superior)
* **Docker** y **Docker Compose**
* **HURL** (para la suite de integración de pruebas HTTP)
* **sqlc** (opcional, para regenerar queries)

## Ejecución Automatizada de Pruebas (Instrucción de Entrega)

Para probar el proyecto completo desde cero (setup de base de datos, ejecución de sqlc, compilación, tests HURL y cleanup final), clonar el repositorio e ingresar a la rama `tp3`:

```bash
git checkout tp3
make test-api
