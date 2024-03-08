CREATE TABLE IF NOT EXISTS characters
(
    ID        bigserial PRIMARY KEY,
    CreatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    UpdatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    FirstName text NOT NULL,
    LastName text NOT NULL,
    House text NOT NULL,
    OriginStatus text NOT NULL
); 