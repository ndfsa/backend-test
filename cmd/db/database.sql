-- function to create a user, checks for null and encrypts the password
DROP FUNCTION IF EXISTS CREATE_USER;
CREATE FUNCTION CREATE_USER(
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


-- function to update user, only username and fullname can be updated
DROP PROCEDURE IF EXISTS UPDATE_USER;
CREATE PROCEDURE UPDATE_USER(
    _id BIGINT,
    _fullname VARCHAR(300),
    _username VARCHAR(100)
) AS
$$
BEGIN
    UPDATE users SET
        fullname = COALESCE(NULLIF(_fullname, ''), fullname),
        username = COALESCE(NULLIF(_username, ''), username)
    WHERE id = _id;
END
$$
LANGUAGE 'plpgsql';


-- authentication function, returns user ID for tokens
DROP FUNCTION IF EXISTS AUTHENTICATE_USER;
CREATE FUNCTION AUTHENTICATE_USER(
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

    SELECT (password = crypt(_password, password)) , id
    INTO _auth, _id
    FROM users
    WHERE username = _username;

    IF _auth THEN
        RETURN _id;
    END IF;

    RAISE EXCEPTION 'user authentication unsuccessful.';
END
$$
LANGUAGE 'plpgsql';
