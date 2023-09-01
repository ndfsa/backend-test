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
    _type SERVICE_TYPE,
    _currency CURR,
    _init_balacne NUMERIC(20, 2)
) RETURNS BIGINT AS
$$
DECLARE
    res BIGINT;
BEGIN
    INSERT INTO services (type, state, currency, init_balance, balance)
    VALUES (_type, 'REQ', _currency, _init_balacne, 0)
    RETURNING id INTO res;

    INSERT INTO user_service (user_id, service_id)
    VALUES (_user_id, res);

    RETURN res;
END
$$
LANGUAGE 'plpgsql';

DROP PROCEDURE IF EXISTS CANCEL_SERVICE;
CREATE PROCEDURE CANCEL_SERVICE(
    _user_id BIGINT,
    _service_id BIGINT
) AS
$$
BEGIN
    UPDATE services
    SET state = 'CLD'
    FROM users JOIN user_service ON users.id = user_id
    WHERE users.id = _user_id AND services.id = _service_id;
END
$$
LANGUAGE 'plpgsql';
