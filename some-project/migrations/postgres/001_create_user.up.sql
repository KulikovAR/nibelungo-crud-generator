-- +migrate Up
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    id INTEGER,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_created_at ON users(created_at);

-- +migrate Down
DROP TABLE users;
