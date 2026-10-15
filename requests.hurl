# 1. Crear una Pizza y capturar el ID asignado por PostgreSQL
POST http://localhost:8080/pizzas
Content-Type: application/json
{
  "nombre": "Muzzarella",
  "descripcion": "Salsa de tomate y queso muzzarella",
  "precio": 1200.50,
  "disponible": true
}

HTTP 201
[Captures]
pizza_id: jsonpath "$.ID"
[Asserts]
jsonpath "$.Nombre" == "Muzzarella"


# 2. Listar Pizzas
GET http://localhost:8080/pizzas

HTTP 200
[Asserts]
jsonpath "$" count >= 1


# 3. Obtener Pizza por el ID dinámico
GET http://localhost:8080/pizzas/{{pizza_id}}

HTTP 200
[Asserts]
jsonpath "$.ID" == {{pizza_id}}
jsonpath "$.Nombre" == "Muzzarella"


# 4. Actualizar Pizza usando el ID dinámico
PUT http://localhost:8080/pizzas/{{pizza_id}}
Content-Type: application/json
{
  "nombre": "Muzzarella Especial",
  "descripcion": "Salsa, muzzarella y aceitunas",
  "precio": 1350.00,
  "disponible": true
}

HTTP 200
[Asserts]
jsonpath "$.Nombre" == "Muzzarella Especial"


# 5. Crear un Pedido asociado a la Pizza recién creada y capturar su ID
POST http://localhost:8080/pedidos
Content-Type: application/json
{
  "cliente_nombre": "Juan Perez",
  "cliente_direccion": "Calle Falsa 123",
  "cliente_telefono": "555-1234",
  "pizza_id": {{pizza_id}},
  "cantidad": 2,
  "total": 2700.00
}

HTTP 201
[Captures]
pedido_id: jsonpath "$.ID"
[Asserts]
jsonpath "$.ClienteNombre" == "Juan Perez"


# 6. Listar Pedidos
GET http://localhost:8080/pedidos

HTTP 200
[Asserts]
jsonpath "$" count >= 1


# 7. Actualizar Estado del Pedido dinámico
PUT http://localhost:8080/pedidos/{{pedido_id}}
Content-Type: application/json
{
  "estado": "en preparacion"
}

HTTP 200
[Asserts]
jsonpath "$.Estado.String" == "en preparacion"


# 8. Eliminar Pedido
DELETE http://localhost:8080/pedidos/{{pedido_id}}

HTTP 204


# 9. Eliminar Pizza
DELETE http://localhost:8080/pizzas/{{pizza_id}}

HTTP 204