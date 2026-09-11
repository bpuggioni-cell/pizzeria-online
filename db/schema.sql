-- Creación de la tabla principal de Pizzas
CREATE TABLE IF NOT EXISTS pizzas (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT NOT NULL,
    precio NUMERIC(10, 2) NOT NULL,
    disponible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pedidos (
    id SERIAL PRIMARY KEY,
    cliente_nombre VARCHAR(100) NOT NULL,
    cliente_direccion TEXT NOT NULL,
    cliente_telefono VARCHAR(50) NOT NULL,
    pizza_id INT NOT NULL,
    cantidad INT NOT NULL DEFAULT 1,
    total NUMERIC(10, 2) NOT NULL,
    estado VARCHAR(50) DEFAULT 'pendiente', -- ej: pendiente, en preparación, entregado
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pizza FOREIGN KEY (pizza_id) REFERENCES pizzas(id) ON DELETE CASCADE
);