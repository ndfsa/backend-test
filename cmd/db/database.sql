DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id BIGSERIAL,
    fullname VARCHAR(300),
    username VARCHAR(100),
    password VARCHAR(60),
    PRIMARY KEY (id),
    UNIQUE (username)
);

INSERT INTO users (fullname, username, password) VALUES (
    'root user',
    'root',
    '$2a$06$DZxsYD5zF5NI/ugKmMmZw.7/hehCmlCpzDOuPutYFmwIlyT37SDGy'
);

CREATE OR REPLACE FUNCTION CREATE_USER(
    _fullname VARCHAR(300),
    _username VARCHAR(100),
    _password VARCHAR(72)
) RETURNS BIGINT AS
$$
DECLARE
    res BIGINT;
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

    INSERT INTO users(fullname, username, password) VALUES (
        _fullname,
        _username,
        crypt(_password, gen_salt('bf')))
    RETURNING id INTO res;

    RETURN res;
END
$$
LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION AUTHENTICATE_USER(
    _username VARCHAR,
    _password VARCHAR
) RETURNS BIGINT AS
$$
DECLARE
    _id BIGINT;
    _auth BOOLEAN;
BEGIN
    IF _username IS NULL OR _username = '' THEN
        RAISE EXCEPTION 'username cannot be NULL or empty.';
    END IF;

    IF _password IS NULL OR _password = '' THEN
        RAISE EXCEPTION 'password cannot be NULL or empty.';
    END IF;

    SELECT (password = crypt(_password, password)) AS pswdmatch, id
    INTO _auth, _id
    FROM users
    WHERE username = _username;

    IF _auth THEN
        RETURN _id;
    END IF;

    RETURN 0;
END
$$
LANGUAGE 'plpgsql';
