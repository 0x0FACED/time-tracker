CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    passport_number VARCHAR(6) UNIQUE NOT NULL,
    pass_serie VARCHAR(4) NOT NULL,
    surname VARCHAR(15),
    name VARCHAR(15),
    patronymic VARCHAR(20),
    address VARCHAR(255)
);