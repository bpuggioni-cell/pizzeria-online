package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq" // Driver de PostgreSQL

	pizzeriaDB "tpmod/internal/db"
	"tpmod/views" // Importá tu paquete de vistas generado por templ
)

type Server struct {
	db *pizzeriaDB.Queries
}

// HANDLER PRINCIPAL (GET /) - Renderiza el SSR completo con las pizzas
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener todas las pizzas de la base de datos usando sqlc
	pizzas, err := s.db.ListPizzas(r.Context())
	if err != nil {
		log.Println("[ERROR handleHome - ListPizzas]:", err)
		http.Error(w, "Error al obtener las pizzas de la base de datos", http.StatusInternalServerError)
		return
	}

	// 2. Renderizar el componente IndexPage pasándole los datos
	component := views.IndexPage(pizzas)
	err = component.Render(r.Context(), w)
	if err != nil {
		log.Println("[ERROR handleHome - Render]:", err)
		http.Error(w, "Error al renderizar la vista", http.StatusInternalServerError)
		return
	}
}

// HANDLER PARA CREAR PIZZA (POST /pizzas) - Tradicional por formulario SSR
func (s *Server) handleCreatePizzaForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")
	descripcion := r.FormValue("descripcion")
	precioStr := r.FormValue("precio")

	if nombre == "" || precioStr == "" {
		http.Error(w, "Nombre y precio son obligatorios", http.StatusBadRequest)
		return
	}

	var precio float64
	_, err := fmt.Sscanf(precioStr, "%f", &precio)
	if err != nil {
		http.Error(w, "Precio inválido", http.StatusBadRequest)
		return
	}

	// Insertar en la base de datos con sqlc
	_, err = s.db.CreatePizza(r.Context(), pizzeriaDB.CreatePizzaParams{
		Nombre:      nombre,
		Descripcion: descripcion,
		Precio:      fmt.Sprintf("%.2f", precio),
		Disponible:  sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		log.Println("[ERROR handleCreatePizzaForm]:", err)
		http.Error(w, "Error al guardar la pizza", http.StatusInternalServerError)
		return
	}

	// Redirigir de vuelta a la página principal para ver la lista actualizada
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST /pizzas/delete/{id} - Eliminar pizza mediante formulario SSR
func (s *Server) handleDeletePizzaForm(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = s.db.DeletePizza(r.Context(), int32(id))
	if err != nil {
		log.Println("[ERROR DeletePizza]:", err)
		http.Error(w, "Error al eliminar la pizza", http.StatusInternalServerError)
		return
	}

	// Redirigir de vuelta a la página principal para ver la lista actualizada
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	// 1. Obtener la cadena de conexión de la base de datos
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/pizzeria?sslmode=disable"
	}

	// 2. Abrir conexión con PostgreSQL
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error al abrir la conexión con la base de datos: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	fmt.Println("Conexión exitosa a la base de datos PostgreSQL")

	// 3. Inicializar queries de sqlc y servidor
	queries := pizzeriaDB.New(dbConn)
	server := &Server{db: queries}

	// 4. Registrar Rutas SSR (Adiós JSON, hola HTML desde Go)
	http.HandleFunc("GET /", server.handleHome)
	http.HandleFunc("POST /pizzas", server.handleCreatePizzaForm)
	http.HandleFunc("POST /pizzas/delete/{id}", server.handleDeletePizzaForm)

	// 5. Iniciar servidor HTTP
	port := ":8080"
	fmt.Printf("Servidor SSR de la Pizzería escuchando en http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}

}
