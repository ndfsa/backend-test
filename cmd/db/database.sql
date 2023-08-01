DROP TABLE IF EXISTS users;
CREATE TABLE users (
    userId BIGSERIAL,
    userFullName VARCHAR(300),
    username VARCHAR(100),
    userPassword VARCHAR(64),
    PRIMARY KEY (userId)
);

INSERT INTO users (userFullName, userName, userPassword) VALUES (
    'root user',
    'root',
    '5503bf09e84c357deb9ac409cd28a1193cf1dffc06c49716d45b87795f700ca1'
);

CREATE OR REPLACE FUNCTION createUser(
    _userFullName VARCHAR(300),
    _userName VARCHAR(100),
    _userPassword VARCHAR(100)
) RETURNS BIGINT AS
$$
DECLARE res BIGINT;
begin

    INSERT INTO users(userFullName, userName, userPassword) VALUES (
        _name,
        _username,
        (SELECT encode(sha256((_username || '|' || _password) ::bytea), 'hex'))
    )
    RETURNING id
    INTO res;

    RETURN res;
end
$$
LANGUAGE 'plpgsql';
