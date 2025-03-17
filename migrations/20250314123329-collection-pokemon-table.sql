-- +migrate Up
CREATE TABLE collection_pokemon (
    pokemon_id uuid NOT NULL,
    user_id uuid NOT NULL,
    caught date,
    nickname text,
    CONSTRAINT fk_pokemon FOREIGN KEY (pokemon_id) REFERENCES pokemon (pokemon_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT collection_pkey PRIMARY KEY (pokemon_id, user_id)
);

CREATE INDEX idx_collection_pokemon_fk_user_id ON collection_pokemon (user_id);

CREATE INDEX idx_collection_pokemon_fk_pokemon_id ON collection_pokemon (pokemon_id);

-- +migrate Down
DROP TABLE IF EXISTS collection_pokemon;

