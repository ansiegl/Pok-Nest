-- +migrate Up
CREATE TYPE pokemon_type AS ENUM (
    'Fire',
    'Water',
    'Grass',
    'Electric',
    'Psychic',
    'Ghost',
    'Dragon',
    'Fairy',
    'Dark',
    'Fighting',
    'Bug',
    'Normal',
    'Rock',
    'Ground',
    'Poison',
    'Steel',
    'Ice'
);

CREATE TABLE pokemon (
    pokemon_id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    "name" text NOT NULL,
    type_1 pokemon_type NOT NULL,
    type_2 pokemon_type,
    generation int NOT NULL,
    legendary boolean NOT NULL,
    CONSTRAINT pokemon_pkey PRIMARY KEY (pokemon_id),
    CONSTRAINT pokemon_name_key UNIQUE ("name")
);

-- +migrate Down
DROP TABLE IF EXISTS pokemon;

DROP TYPE IF EXISTS pokemon_type;

