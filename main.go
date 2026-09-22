package main

import (
	"fmt"
	"net/http"
)

// Handler para la ruta inicial
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}

	// Servir el archivo físico index.html
	http.ServeFile(w, r, "index.html")
}

// Handler para rutas inexistentes (Error 404)
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404 - Página no encontrada</h1><a href='/'>Volver al inicio</a>")
}

func main() {
	http.HandleFunc("/", homeHandler)

	port := ":8080"
	fmt.Printf("Servidor de la Pizzería escuchando en http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
