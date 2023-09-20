CREATE TABLE users
(
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(255) not null,
    surname     VARCHAR(255) not null,
    patronymic  VARCHAR(255) default null,
    age         INT          default -1,
    gender      VARCHAR(255) default null,
    nationality VARCHAR(255) default null
)