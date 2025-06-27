-- +migrate Up
CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    id INTEGER,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    price FLOAT NOT NULL,
    created_at VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_created_at ON products(created_at);

-- +migrate Down
DROP TABLE products;
