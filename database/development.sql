CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    username VARCHAR(128) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    session VARCHAR(36),
    expiry TIMESTAMP
)
