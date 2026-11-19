# Sistema de Gestión de Pizzería - SSR (Trabajo Práctico 5)

Aplicación web desarrollada en **Go** utilizando **Server-Side Rendering (SSR)** con componentes **`templ`**, estilos basados en **Pico.css**, persistencia de datos mediante **PostgreSQL** y mapeo tipado con **`sqlc`**. Implementa estrictamente el patrón **Post/Redirect/Get (PRG)**.

## Requisitos Previos
* **Go** (versión 1.22 o superior)
* **Docker** (daemon en ejecución, utilizado para levantar PostgreSQL de forma aislada)
* **GNU Make**
* Herramienta **`templ`** instalada (`go install [github.com/a-h/templ/cmd/templ@latest](https://github.com/a-h/templ/cmd/templ@latest)`)
* Herramienta **`sqlc`** instalada (`go install [github.com/sqlc-dev/sqlc/cmd/sqlc@latest](https://github.com/sqlc-dev/sqlc/cmd/sqlc@latest)`)

## Instrucciones de Ejecución

Para preparar el entorno de base de datos, aplicar el esquema, compilar y ejecutar la aplicación de forma automatizada mediante el `Makefile`, ejecute en la raíz del proyecto:

## Cambios Realizados en la Reentrega
* **Carga automática del esquema:** Se incorporó la regla `schema-init` en `make run-app` para ejecutar `db/schema.sql` sobre la base `pizzeria` antes de iniciar el servidor, evitando errores por tablas inexistentes al arrancar desde una base limpia.
* **Limpieza aislada al proyecto:** Se reemplazó `docker volume prune -f` global por un volumen nombrado propio del proyecto (`pizzeria-pgdata`), eliminando únicamente el contenedor `pizzeria-postgres` y dicho volumen.
* **Manejador de salida (`trap`):** Se agregó un `trap` en `make run-app` sobre `EXIT`, `INT` y `TERM` para garantizar la limpieza automática de los recursos del proyecto al finalizar o interrumpir la ejecución (`Ctrl+C`).

```bash
make run-app
```