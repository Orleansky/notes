DROP TABLE IF EXISTS notes, users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
	password TEXT NOT NULL
);

CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL
);