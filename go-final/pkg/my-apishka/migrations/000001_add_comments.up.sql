CREATE TABLE IF NOT EXISTS comments (
    Id SERIAL PRIMARY KEY,
    UsernameID bigserial NOT NULL REFERENCES users ON DELETE CASCADE,
    Comment TEXT NOT NULL,
    CharacterID bigserial NOT NULL REFERENCES characters ON DELETE CASCADE
);