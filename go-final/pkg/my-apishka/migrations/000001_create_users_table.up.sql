CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users 
(
    ID bigserial PRIMARY KEY,
    CreatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW() ,
    Username text NOT NULL,
    Email citext UNIQUE NOT NULL,
    Password bytea NOT NULL,
    Activated bool NOT NULL,
    Version integer NOT NULL DEFAULT 1
);