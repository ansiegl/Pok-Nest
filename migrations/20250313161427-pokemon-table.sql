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
    'Flying',
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
    hp integer NOT NULL,
    attack integer NOT NULL,
    defense integer NOT NULL,
    speed integer NOT NULL,
    special integer NOT NULL,
    gif_url text NOT NULL,
    png_url text NOT NULL,
    "description" text NOT NULL,
    CONSTRAINT pokemon_pkey PRIMARY KEY (pokemon_id),
    CONSTRAINT pokemon_name_key UNIQUE ("name")
);

-- +migrate Down
DROP TABLE IF EXISTS pokemon;

DROP TYPE IF EXISTS pokemon_type;

