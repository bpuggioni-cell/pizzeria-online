-- name: CreatePizza :one
INSERT INTO pizzas (nombre, descripcion, precio, disponible) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: GetPizza :one
SELECT * FROM pizzas 
WHERE id = $1;

-- name: ListPizzas :many
SELECT * FROM pizzas 
ORDER BY id ASC;

-- name: UpdatePizza :one
UPDATE pizzas 
SET nombre = $2, descripcion = $3, precio = $4, disponible = $5 
WHERE id = $1
RETURNING *;

-- name: DeletePizza :exec
DELETE FROM pizzas 
WHERE id = $1;



-- name: CreatePedido :one
INSERT INTO pedidos (cliente_nombre, cliente_direccion, cliente_telefono, pizza_id, cantidad, total, estado) 
VALUES ($1, $2, $3, $4, $5, $6, $7) 
RETURNING *;

-- name: GetPedido :one
SELECT * FROM pedidos 
WHERE id = $1;

-- name: ListPedidos :many
SELECT * FROM pedidos 
ORDER BY created_at DESC;

-- name: UpdatePedidoState :one
UPDATE pedidos
SET estado = $2
WHERE id = $1
RETURNING *;

-- name: DeletePedido :exec
DELETE FROM pedidos 
WHERE id = $1;