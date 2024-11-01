BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS users(
    id       INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login    VARCHAR(50) NOT NULL,
    password VARCHAR(250) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS login_idx on users (login);


CREATE TABLE IF NOT EXISTS user_records(
    id         INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name       VARCHAR(50) NOT NULL,
    data       BYTEA       NOT NULL,
    data_type  VARCHAR(50) NOT NULL CHECK ( data_type IN ('password', 'card', 'text', 'file')),
    version    INT         NOT NULL DEFAULT 1,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id    INT         NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_id_name_idx on user_records (user_id, name);

COMMIT;