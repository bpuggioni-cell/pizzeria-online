# Sistema de Gestión de Pizzería - Interfaces Reactivas con HTMX (TP6)

Aplicación web evolucionada para comportarse como una **Single Page Application (SPA)** mediante **Server-Side Rendering (SSR)** utilizando **Go**, componentes **`templ`**, estilos basados en **Pico.css**, persistencia con **PostgreSQL**, mapeo tipado con **`sqlc`** y dinamismo interactivo con **HTMX** (sin recargas de página).

## Características Principales
* **Cero JavaScript nativo:** Toda la interactividad asíncrona se maneja a través de atributos declarativos de HTMX en los componentes HTML (`hx-post`, `hx-delete`, `hx-target`, `hx-swap`).
* **Renderizado del lado del servidor (SSR):** Los fragmentos visuales son generados directamente por componentes `templ` en Go.
* **Ciclo de vida automatizado:** Gestión completa de contenedores Docker, volúmenes, espera de disponibilidad de la base de datos y compilación mediante un `Makefile`.

## Requisitos del Sistema
* **Go** (versión 1.22 o superior)
* **Docker y Docker Compose** (para el motor de base de datos)
* Herramienta **`templ`** instalada en el sistema (`go install github.com/a-h/templ/cmd/templ@latest`)

## Instrucciones de Ejecución

Para clonar, preparar el entorno de base de datos, compilar y ejecutar la aplicación de forma totalmente automatizada, ejecute el siguiente comando en la raíz del proyecto:

```bash
make run-app