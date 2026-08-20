package main

import (
	"fmt"
	"net/http"
)

const indexHTML = `<!DOCTYPE html>
	<html lang="es">
	<head>
		<title>Pizzería Online</title>
	</head>
	<body>
    <h1>Sistema de Pedidos de Pizzería</h1>
    <p>
        Bienvenido al sistema de gestión de pedidos. Esta aplicación te permite realizar 
        pedidos de pizza, consultar el estado de tu entrega y gestionar la lista de pedidos en tiempo real.
    </p>
	</body>
	</html>`

// Handler para la ruta inicial
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Validar que la ruta sea exactamente "/" y no otra inexistente como /abc
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, indexHTML)
}

// Handler para rutas inexistentes (Error 404)
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404 - Página no encontrada</h1><a href='/'>Volver al inicio</a>")
}

func main() {

	http.HandleFunc("/", homeHandler)
	// 3. Arrancar servidor en puerto 8080
	port := ":8080"
	fmt.Printf("Servidor de la Pizzería escuchando en http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
