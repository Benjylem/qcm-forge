-- Une session = un jeton aléatoire opaque, stocké hashé (comme un mot
-- de passe) pour qu'un vol de la base ne donne pas directement des
-- sessions valides. ON DELETE CASCADE : supprimer un user supprime
-- ses sessions, pas de ligne orpheline.
CREATE TABLE sessions (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    token_hash bytea       NOT NULL UNIQUE,
    user_id    bigint      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- Chaque requête authentifiée cherche par token_hash : index nécessaire.
CREATE INDEX sessions_token_hash_idx ON sessions (token_hash);
