CREATE TABLE IF NOT EXISTS users(
    id bigserial primary key,
    email varchar (300) not null unique,
    password varchar (100) not null,
    active boolean,
    email_verified boolean,
    email_subscribed boolean,
    admin boolean,
    created TIMESTAMPTZ,
    last_login TIMESTAMPTZ,
    last_action TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS refresh_token (
    id uuid PRIMARY KEY,
    token_string varchar (215),
    expires BIGINT
)