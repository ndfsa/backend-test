DROP FUNCTION IF EXISTS GET_USER_SERVICES;
CREATE FUNCTION GET_USER_SERVICES(
    _id BIGINT
) RETURNS SETOF SERVICES AS
$$
BEGIN
    RETURN QUERY SELECT s.id, s.type, s.state, s.currency, s.init_balance, s.balance
    FROM users u
    JOIN user_service us ON u.id = us.user_id
    JOIN services s ON s.id = us.service_id
    WHERE u.id = _id;
END
$$
LANGUAGE 'plpgsql';


DROP FUNCTION IF EXISTS CREATE_SERVICE;
CREATE FUNCTION CREATE_SERVICE(
    _user_id BIGINT,
    _type SMALLINT,
    _currency CURR,
    _init_balacne NUMERIC(20, 2)
) RETURNS BIGINT AS
$$
DECLARE
    res BIGINT;
BEGIN
    INSERT INTO services (type, state, currency, init_balance, balance)
    VALUES (_type, 1, _currency, _init_balacne, 0)
    RETURNING id INTO res;

    INSERT INTO user_service (user_id, service_id)
    VALUES (_user_id, res);

    RETURN res;
END
$$
LANGUAGE 'plpgsql';
