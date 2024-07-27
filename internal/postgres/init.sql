CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    author TEXT NOT NULL,
    text TEXT NOT NULL,
    created TIMESTAMP,
    sent BOOLEAN NOT NULL
);