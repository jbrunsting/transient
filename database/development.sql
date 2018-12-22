CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    username VARCHAR(128) NOT NULL,
    password TEXT NOT NULL,
    session VARCHAR(36),
    expiry VARCHAR(23)
)
