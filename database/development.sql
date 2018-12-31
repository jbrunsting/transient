CREATE TABLE IF NOT EXISTS Users (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    username VARCHAR(128) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Sessions (
    id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sessionId VARCHAR(36) NOT NULL PRIMARY KEY,
    expiry TIMESTAMP NOT NULL
);
