CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    passport_number VARCHAR(6) UNIQUE NOT NULL,
    pass_serie VARCHAR(4) NOT NULL,
    surname VARCHAR(15),
    name VARCHAR(15),
    patronymic VARCHAR(20),
    address VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    description VARCHAR(1023),
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
