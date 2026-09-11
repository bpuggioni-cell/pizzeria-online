Pizzería - TP2: Capa de Persistencia

Este repositorio contiene la implementación de la capa de persistencia para la aplicación de la Pizzería, utilizando PostgreSQL y sqlc para la generación de código Go fuertemente tipado.

Estructura y Arquitectura

- Base de Datos:
  -`pizzas`: Menú con detalles de productos, precios y disponibilidad.
  -`pedidos`: Almacena la orden solicitada junto a los datos de entrega del cliente (`cliente_nombre`, `cliente_direccion`, `cliente_telefono`), garantizando el flujo de pedidos sin necesidad de registro previo.

- Archivos Clave:
  - `schema.sql`: Sentencias DDL para la creación de tablas.
  - `queries.sql`: Consultas SQL con anotaciones para `sqlc`.
  - `sqlc.yaml`: Configuración de generación de código Go.
  - `internal/db/`: Código Go autogenerado (`models.go`, `db.go`, `queries.sql.go`) y tests de integración (`db_test.go`).
  - `test.sh` y `Makefile`: Automatización de entorno en Docker y ejecución de tests.

## Ejecución de Tests

Para ejecutar el ciclo de pruebas automatizado (levantar contenedor PostgreSQL, aplicar schema, generar código, ejecutar `go test` y limpiar el entorno):

```bash
make test
