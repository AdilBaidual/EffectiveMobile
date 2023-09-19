CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) not null,
    surname    VARCHAR(255) not null,
    patronymic VARCHAR(255) default null,
    age        INT          default null,
    gender     VARCHAR(255) default null,
    nationality VARCHAR(255) default null
)