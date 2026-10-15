package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq" // Driver de PostgreSQL

	pizzeriaDB "tpmod/internal/db"
)

type Server struct {
	db *pizzeriaDB.Queries
}

// Structs DTO para decodificar JSON entrante correctamente
type PizzaRequest struct {
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	Disponible  bool    `json:"disponible"`
}

type PedidoRequest struct {
	ClienteNombre    string  `json:"cliente_nombre"`
	ClienteDireccion string  `json:"cliente_direccion"`
	ClienteTelefono  string  `json:"cliente_telefono"`
	PizzaID          int32   `json:"pizza_id"`
	Cantidad         int32   `json:"cantidad"`
	Total            float64 `json:"total"`
}

type UpdateEstadoPedidoRequest struct {
	Estado string `json:"estado"`
}

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

// HANDLERS PARA PIZZAS

// POST /pizzas - Crear una pizza
func (s *Server) handleCreatePizza(w http.ResponseWriter, r *http.Request) {
	var req PizzaRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[ERROR handleCreatePizza]:", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" {
		http.Error(w, "El nombre de la pizza es obligatorio", http.StatusBadRequest)
		return
	}

	nuevaPizza, err := s.db.CreatePizza(r.Context(), pizzeriaDB.CreatePizzaParams{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Precio:      fmt.Sprintf("%.2f", req.Precio),
		Disponible:  sql.NullBool{Bool: req.Disponible, Valid: true},
	})
	if err != nil {
		http.Error(w, "Error al crear la pizza en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevaPizza)
}

// GET /pizzas - Listar Pizzas
func (s *Server) handleListPizzas(w http.ResponseWriter, r *http.Request) {
	pizzas, err := s.db.ListPizzas(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener las pizzas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizzas)
}

// GET /pizzas/{id} - Obtener Pizza por ID
func (s *Server) handleGetPizzaByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	pizza, err := s.db.GetPizza(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Pizza no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizza)
}

// PUT /pizzas/{id} - Actualizar pizza por ID
func (s *Server) handleUpdatePizza(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "El ID provisto es inválido", http.StatusBadRequest)
		return
	}

	var req PizzaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" {
		http.Error(w, "El nombre de la pizza no puede estar vacío", http.StatusBadRequest)
		return
	}

	pizzaActualizada, err := s.db.UpdatePizza(r.Context(), pizzeriaDB.UpdatePizzaParams{
		ID:          int32(id),
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Precio:      fmt.Sprintf("%.2f", req.Precio),
		Disponible:  sql.NullBool{Bool: req.Disponible, Valid: true},
	})
	if err != nil {
		http.Error(w, "Error al actualizar la pizza", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizzaActualizada)
}

// DELETE /pizzas/{id} - Eliminar pizza por ID
func (s *Server) handleDeletePizza(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "El ID provisto es inválido", http.StatusBadRequest)
		return
	}

	err = s.db.DeletePizza(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar la pizza", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// HANDLERS PARA PEDIDOS

// POST /pedidos - Crear un pedido
func (s *Server) handleCreatePedido(w http.ResponseWriter, r *http.Request) {
	var req PedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[ERROR handleCreatePedido]:", err) // <--- Agrega este log
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validaciones básicas del cliente y pedido
	if req.ClienteNombre == "" || req.ClienteDireccion == "" || req.ClienteTelefono == "" {
		http.Error(w, "Los datos del cliente son obligatorios", http.StatusBadRequest)
		return
	}
	if req.PizzaID <= 0 {
		http.Error(w, "Debe especificar una pizza válida (pizza_id)", http.StatusBadRequest)
		return
	}

	nuevoPedido, err := s.db.CreatePedido(r.Context(), pizzeriaDB.CreatePedidoParams{
		ClienteNombre:    req.ClienteNombre,
		ClienteDireccion: req.ClienteDireccion,
		ClienteTelefono:  req.ClienteTelefono,
		PizzaID:          req.PizzaID,
		Cantidad:         req.Cantidad,
		Total:            fmt.Sprintf("%.2f", req.Total),
	})
	if err != nil {
		http.Error(w, "Error al registrar el pedido", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(nuevoPedido)
}

// GET /pedidos - Listar todos los pedidos
func (s *Server) handleListPedidos(w http.ResponseWriter, r *http.Request) {
	pedidos, err := s.db.ListPedidos(r.Context())
	if err != nil {
		http.Error(w, "Error al consultar la lista de pedidos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidos)
}

// GET /pedidos/{id} - Obtener pedido por ID
func (s *Server) handleGetPedidoByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	pedido, err := s.db.GetPedido(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Pedido no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedido)
}

// PUT /pedidos/{id} - Actualizar el estado de un pedido
func (s *Server) handleUpdatePedidoState(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "El ID provisto es inválido", http.StatusBadRequest)
		return
	}

	var req UpdateEstadoPedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	pedidoActualizado, err := s.db.UpdatePedidoState(r.Context(), pizzeriaDB.UpdatePedidoStateParams{
		ID:     int32(id),
		Estado: sql.NullString{String: req.Estado, Valid: true},
	})
	if err != nil {
		http.Error(w, "Error al actualizar el estado del pedido o pedido no encontrado", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidoActualizado)
}

// DELETE /pedidos/{id} - Eliminar/Cancelar un pedido por ID
func (s *Server) handleDeletePedido(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = s.db.DeletePedido(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar el pedido", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// 1. Obtener la URL/cadena de conexión de la base de datos desde las variables de entorno
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Valor por defecto para desarrollo local si la variable no está seteada
		dbURL = "postgres://postgres:postgres@localhost:5432/pizzeria?sslmode=disable"
	}

	// 2. Abrir la conexión con PostgreSQL
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error al abrir la conexión con la base de datos: %v", err)
	}
	defer dbConn.Close()

	// 3. Verificar que la base de datos esté accesible
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	fmt.Println("Conexión exitosa a la base de datos PostgreSQL")

	// 4. Crear la instancia del repositorio generado por sqlc
	queries := pizzeriaDB.New(dbConn)

	// 5. Instanciar el servidor con el repositorio
	server := &Server{db: queries}

	// Handler para la raíz / estáticos
	http.HandleFunc("/", homeHandler)

	// Endpoints API REST - PIZZAS
	http.HandleFunc("POST /pizzas", server.handleCreatePizza)
	http.HandleFunc("GET /pizzas", server.handleListPizzas)
	http.HandleFunc("GET /pizzas/{id}", server.handleGetPizzaByID)
	http.HandleFunc("PUT /pizzas/{id}", server.handleUpdatePizza)
	http.HandleFunc("DELETE /pizzas/{id}", server.handleDeletePizza)

	// Endpoints API REST - PEDIDOS
	http.HandleFunc("POST /pedidos", server.handleCreatePedido)
	http.HandleFunc("GET /pedidos", server.handleListPedidos)
	http.HandleFunc("GET /pedidos/{id}", server.handleGetPedidoByID)
	http.HandleFunc("PUT /pedidos/{id}", server.handleUpdatePedidoState)
	http.HandleFunc("DELETE /pedidos/{id}", server.handleDeletePedido)

	// 6. Iniciar el servidor HTTP
	port := ":8080"
	fmt.Printf("Servidor de la Pizzería escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
