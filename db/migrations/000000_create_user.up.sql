CREATE TABLE "user"
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    email   VARCHAR(255) NOT NULL
);