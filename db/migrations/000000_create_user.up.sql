CREATE TABLE "user"
(
    id      SERIAL,
    name    VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    email   VARCHAR(255) NOT NULL,
    password_hash VARCHAR(128) NOT NULL,
    salt_hash VARCHAR(128) NOT NULL,
    enabled BOOLEAN NOT NULL,

    CONSTRAINT user_pk PRIMARY KEY (id),
    CONSTRAINT user_ak01 UNIQUE (email)
);