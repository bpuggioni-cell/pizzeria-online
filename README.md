# Sistema de Gestión de Pizzería - SSR (Trabajo Práctico 5)

Aplicación web desarrollada en **Go** utilizando **Server-Side Rendering (SSR)** con componentes **`templ`**, estilos basados en **Pico.css**, persistencia de datos mediante **PostgreSQL** y mapeo tipado con **`sqlc`**. Implementa estrictamente el patrón **Post/Redirect/Get (PRG)**.

## Requisitos Previos
* **Go** (versión 1.22 o superior)
* **Docker y Docker Compose** (para el manejo automatizado de la base de datos)
* Herramienta **`templ`** instalada en el sistema (`go install github.com/a-h/templ/cmd/templ@latest`)

## Instrucciones de Ejecución

Para clonar, preparar el entorno de base de datos, compilar y ejecutar la aplicación de forma totalmente automatizada mediante el Makefile, ejecute el siguiente comando en la raíz del proyecto:

```bash
make run-app