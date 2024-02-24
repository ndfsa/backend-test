DROP TABLE IF EXISTS user_service CASCADE;
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id UUID,
    fullname VARCHAR(300) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL,
    PRIMARY KEY (id)
);

DROP TYPE IF EXISTS CURRENCY CASCADE;
-- USD United States Dollar
-- CAD Canadian Dollar
-- JPY Japanese Yen
-- NOK Norwegian crown
CREATE TYPE CURRENCY AS ENUM ('USD', 'CAD', 'JPY', 'NOK');

DROP TYPE IF EXISTS SERVICE_TYPE CASCADE;
-- SVA Savings account
-- CQA Checking account
-- LON Loan
CREATE TYPE SERVICE_TYPE AS ENUM ('SVA', 'CQA', 'LON');

DROP TYPE IF EXISTS SERVICE_STATE CASCADE;
-- REQ Requested service
-- ACT Active service
-- FRZ Frozen service
-- CLD Cancelled service
CREATE TYPE SERVICE_STATE AS ENUM ('REQ', 'ACT', 'FRZ', 'CLD');

DROP TABLE IF EXISTS services CASCADE;
CREATE TABLE services (
    id UUID,
    type SERVICE_TYPE,
    state SERVICE_STATE,
    currency CURRENCY,
    init_balance NUMERIC(20, 2),
    balance NUMERIC(20, 2),
    PRIMARY KEY (id)
);

DROP TYPE IF EXISTS TRANSACTION_STATE CASCADE;
-- INITIAL Initial state
-- PROCESSING Processing transaction
-- INVESTIGATION Under investigation
-- ERROR Error processing
-- DONE Transaction processed
-- CANCELLED Cancelled transaction
-- ROLLBACK Rolled back transaction
CREATE TYPE TRANSACTION_STATE AS ENUM (
    'INITIAL',
    'PROCESSING',
    'INVESTIGATION',
    'ERROR',
    'DONE',
    'CANCELLED',
    'ROLLBACK'
);

DROP TABLE IF EXISTS transactions CASCADE;
CREATE TABLE transactions (
    id UUID,
    state SMALLINT,
    currency CURRENCY,
    amount NUMERIC(20, 2),
    from UUID,
    to UUID,
    PRIMARY KEY (id),
    FOREIGN KEY (from) REFERENCES services ON DELETE CASCADE,
    FOREIGN KEY (to) REFERENCES services ON DELETE CASCADE,
);

CREATE TABLE user_service (
    user_id UUID,
    service_id UUID,
    PRIMARY KEY (user_id, service_id),
    FOREIGN KEY (user_id) REFERENCES users ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES services ON DELETE CASCADE
);

CREATE USER back WITH PASSWORD 'root';
GRANT ALL PRIVILEGES ON DATABASE cardboard_bank TO back;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO back;
