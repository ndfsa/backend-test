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
