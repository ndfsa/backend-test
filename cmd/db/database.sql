DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id BIGSERIAL,
    fullname VARCHAR(300),
    username VARCHAR(100),
    password VARCHAR(64),
    passsalt VARCHAR(64),
    PRIMARY KEY (id),
    UNIQUE (username)
);

INSERT INTO users (fullname, username, password, passsalt) VALUES (
    'root user',
    'root',
    '48c5364e00864c5366d5a220559c04be5376af5ef5dfdf1e5008731334560fd4',
    '5503bf09e84c357deb9ac409cd28a1193cf1dffc06c49716d45b87795f700ca1'
);

CREATE OR REPLACE FUNCTION CREATE_USER(
    _fullname VARCHAR(300),
    _username VARCHAR(100),
    _password VARCHAR(100)
) RETURNS BIGINT AS
$$
DECLARE
    res BIGINT;
    password_salt VARCHAR(64);
BEGIN
    IF _fullname IS NULL OR _fullname = '' THEN
        RAISE EXCEPTION 'fullname cannot be NULL or empty.';
    END IF;

    IF _username IS NULL OR _username = '' THEN
        RAISE EXCEPTION 'username cannot be NULL or empty.';
    END IF;

    IF _password IS NULL OR _password = '' THEN
        RAISE EXCEPTION 'password cannot be NULL or empty.';
    END IF;

    password_salt := gen_salt('sha256');
    INSERT INTO users(fullname, username, password, passsalt) VALUES (
        _fullname,
        _username,
        encode(sha256((password_salt || '|' || _password) ::bytea), 'hex'),
        password_salt
    ) RETURNING id INTO res;

    RETURN res;
end
$$
LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION AUTHENTICATE_USER(
    _username VARCHAR(100),
    _password VARCHAR(100)
) RETURNS BIGINT AS
$$
DECLARE
    res BIGINT;
    _stored_id BIGINT;
    _stored_password VARCHAR(64);
    _stored_passsalt VARCHAR(64);
    _hashed_password VARCHAR(64);
BEGIN
    IF _username IS NULL OR _username = '' THEN
        RAISE EXCEPTION 'username cannot be NULL or empty.';
    END IF;

    IF _password IS NULL OR _password = '' THEN
        RAISE EXCEPTION 'password cannot be NULL or empty.';
    END IF;

    SELECT id, password, passsalt
    INTO _stored_id, _stored_password, _stored_passsalt
    FROM users
    WHERE username = _username;

    RAISE NOTICE 'user %', _stored_id;
    RAISE NOTICE 'pass %', _stored_password;
    RAISE NOTICE 'salt %', _stored_passsalt;

    _hashed_password := encode(sha256((_stored_passsalt || '|' || _password) ::bytea), 'hex');

    RAISE NOTICE 'hashed %', _hashed_password;
    IF _stored_password <> _hashed_password THEN
        RAISE EXCEPTION 'authentication unsuccessful';
    END IF;

    RETURN _stored_id;
end
$$
LANGUAGE 'plpgsql';
