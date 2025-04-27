CREATE TABLE IF NOT EXISTS listings (
    id        INTEGER PRIMARY KEY,
    title     TEXT NOT NULL,
    description TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    category TEXT NOT NULL,
    closed BOOLEAN NOT NULL,
    price INTEGER NOT NULL,
    creator INTEGER NOT NULL
);