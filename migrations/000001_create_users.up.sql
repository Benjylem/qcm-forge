-- Comptes utilisateurs. Le mot de passe n'est jamais stocké : seulement
-- son hash (argon2id/bcrypt, calculé par l'API).
CREATE TABLE users (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email         text        NOT NULL,
    password_hash text        NOT NULL,
    role          text        NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at    timestamptz NOT NULL DEFAULT now()
);

-- Unicité insensible à la casse : Ben@x.fr et ben@x.fr sont le même compte.
CREATE UNIQUE INDEX users_email_lower_idx ON users (lower(email));
