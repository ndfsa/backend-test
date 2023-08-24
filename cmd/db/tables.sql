DROP TABLE IF EXISTS user_service CASCADE;
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id BIGSERIAL,
    fullname VARCHAR(300) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL,
    PRIMARY KEY (id)
);

DROP TYPE CURR CASCADE;
CREATE TYPE curr AS ENUM ('USD', 'CAD', 'JPY', 'NOK');

DROP TABLE IF EXISTS services CASCADE;
CREATE TABLE services (
    id BIGSERIAL,
    type SMALLINT,
    state SMALLINT,
    currency CURR,
    init_balance NUMERIC(20, 2),
    balance NUMERIC(20, 2),
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS transactions CASCADE;
CREATE TABLE transactions (
    id BIGSERIAL,
    type SMALLINT,
    state SMALLINT,
    currency CURR,
    amount NUMERIC(20, 2),
    PRIMARY KEY (id)
);

CREATE TABLE user_service (
    user_id BIGINT,
    service_id BIGINT,
    PRIMARY KEY (user_id, service_id),
    FOREIGN KEY (user_id) REFERENCES users ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES services ON DELETE CASCADE
);

-- transactions from world or to world, are represented as from NULL or to NULL respectively
CREATE TABLE service_transaction (
    from_service_id BIGINT,
    to_service_id BIGINT,
    transaction_id BIGINT,
    FOREIGN KEY (from_service_id) REFERENCES services ON DELETE CASCADE,
    FOREIGN KEY (to_service_id) REFERENCES services ON DELETE CASCADE,
    FOREIGN KEY (transaction_id) REFERENCES users ON DELETE CASCADE
);

-- create root user with default password
INSERT INTO users (fullname, username, password) VALUES (
    'root user',
    'root',
    '$2a$06$DZxsYD5zF5NI/ugKmMmZw.7/hehCmlCpzDOuPutYFmwIlyT37SDGy'
);
