CREATE TABLE IF NOT EXISTS genres (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT genres_name_uq UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS genres_name_idx ON genres(name);
