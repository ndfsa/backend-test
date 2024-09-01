DROP TYPE IF EXISTS USER_ROLE CASCADE;
-- USR Regular
-- OFC Officer
-- ADM Administrator
CREATE TYPE USER_ROLE AS ENUM ('USR', 'OFC', 'ADM');

DROP TABLE IF EXISTS user_service CASCADE;
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id UUID,
    role USER_ROLE,
    fullname VARCHAR(300) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL,
    PRIMARY KEY (id)
);

DROP TYPE IF EXISTS CURRENCY CASCADE;
-- USD United States Dollar
-- CAD Canadian Dollar
-- JPY Japanese Yen
-- NOK Norwegian Crown
CREATE TYPE CURRENCY AS ENUM ('USD', 'CAD', 'JPY', 'NOK');

DROP TYPE IF EXISTS SERVICE_TYPE CASCADE;
-- SAV Savings
-- CHQ Chequing
-- LOA Loan
-- LOC Line of credit
-- COD Certificate of deposit
CREATE TYPE SERVICE_TYPE AS ENUM ('SAV', 'CHQ', 'LOA', 'LOC', 'COD');

DROP TYPE IF EXISTS SERVICE_STATE CASCADE;
-- REQ Requested service
-- ACT Active service
-- FRZ Frozen service
-- CLS Closed service
CREATE TYPE SERVICE_STATE AS ENUM ('REQ', 'ACT', 'FRZ', 'CLS');

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
-- PRC Processing
-- ERR Error
-- SUC Success
CREATE TYPE TRANSACTION_STATE AS ENUM ('PRC', 'ERR', 'SUC');

DROP TABLE IF EXISTS transactions CASCADE;
CREATE TABLE transactions (
    id UUID,
    state TRANSACTION_STATE,
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    currency CURRENCY,
    amount NUMERIC(20, 2),
    source UUID,
    destination UUID,
    PRIMARY KEY (id),
    FOREIGN KEY (source) REFERENCES services ON DELETE CASCADE,
    FOREIGN KEY (destination) REFERENCES services ON DELETE CASCADE
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
