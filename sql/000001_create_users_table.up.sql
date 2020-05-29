CREATE TABLE IF NOT EXISTS users(
    user_id serial PRIMARY KEY,
    email VARCHAR (300) UNIQUE NOT NULL,
    password_hash VARCHAR (300) NOT NULL
);
