DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id BIGSERIAL,
    name VARCHAR(300),
    username VARCHAR(100),
    pass_hash VARCHAR(64),
    PRIMARY KEY (id)
);

INSERT INTO users (name, username, pass_hash) VALUES (
    'root user',
    'root',
    '5503bf09e84c357deb9ac409cd28a1193cf1dffc06c49716d45b87795f700ca1'
);

CREATE OR REPLACE FUNCTION create_user(
    _name VARCHAR(300),
    _username VARCHAR(100),
    _password VARCHAR(100)
) RETURNS BIGINT AS
$$
DECLARE res BIGINT;
begin

    INSERT INTO users(name, username, pass_hash) VALUES (
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
