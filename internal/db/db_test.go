package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	db "tpmod/internal/db"

	_ "github.com/lib/pq"
)

func TestPizzasCRUD(t *testing.T) {
	// 1. Obtener cadena de conexión (o usar valor por defecto para tests)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:1234@localhost:5432/pizzeria?sslmode=disable"
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		t.Fatalf("No se pudo hacer ping a la base de datos: %v", err)
	}

	queries := db.New(conn)
	ctx := context.Background()

	// 2. Probar CREATE
	newPizza, err := queries.CreatePizza(ctx, db.CreatePizzaParams{
		Nombre:      "Muzzarella Test",
		Descripcion: "Salsa de tomate, muzzarella y orégano",
		Precio:      "1200.00",
		Disponible:  sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		t.Fatalf("Error al crear la pizza: %v", err)
	}

	if newPizza.ID == 0 {
		t.Errorf("Se esperaba un ID de pizza válido, pero se obtuvo 0")
	}

	// 3. Probar GET por ID
	fetchedPizza, err := queries.GetPizza(ctx, newPizza.ID)
	if err != nil {
		t.Fatalf("Error al obtener la pizza por ID: %v", err)
	}

	if fetchedPizza.Nombre != newPizza.Nombre {
		t.Errorf("Se esperaba la pizza '%s', pero se obtuvo '%s'", newPizza.Nombre, fetchedPizza.Nombre)
	}

	// 4. Probar LIST
	pizzas, err := queries.ListPizzas(ctx)
	if err != nil {
		t.Fatalf("Error al listar las pizzas: %v", err)
	}

	if len(pizzas) == 0 {
		t.Errorf("Se esperaba al menos una pizza en la lista")
	}

	// 5. Probar DELETE
	err = queries.DeletePizza(ctx, newPizza.ID)
	if err != nil {
		t.Fatalf("Error al eliminar la pizza: %v", err)
	}

	fmt.Println("--> [TEST OK] CRUD de Pizzas con sqlc verificado correctamente.")
}
