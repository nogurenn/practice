CREATE TABLE IF NOT EXISTS artists(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    spotify_artist_id text NOT NULL,
    spotify_uri text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT artists_spotify_artist_id_uq UNIQUE (spotify_artist_id),
    CONSTRAINT artists_spotify_uri_uq UNIQUE (spotify_uri)
);

-- Create a case-insensitive index on the artist name
-- text_pattern_ops: allows the case-insensitive index to work with LIKE 'xxx%'
CREATE INDEX IF NOT EXISTS artists_name_lower_idx ON artists(lower(name) text_pattern_ops);
/*
    From https://www.postgresql.org/docs/current/indexes-opclass.html
    
    Note that you should also create an index with the default operator class if you want queries
    involving ordinary <, <=, >, or >= comparisons to use an index. Such queries cannot use the
    xxx_pattern_ops operator classes.
*/
CREATE INDEX IF NOT EXISTS artists_name_idx ON artists(name);

CREATE INDEX IF NOT EXISTS artists_spotify_artist_id_idx ON artists(spotify_artist_id);

CREATE INDEX IF NOT EXISTS artists_spotify_uri_idx ON artists(spotify_uri);

CREATE TABLE IF NOT EXISTS artist_genres(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    artist_id uuid NOT NULL,
    genre_id uuid NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT artist_genres_artist_id_genre_id_uq UNIQUE (artist_id, genre_id),
    /*
        Cascade-deleting an artist should also delete all of its genres,
        but cascade-deleting a genre with associated artists seems like a business decision that could go either way.
        For now, we'll just cascade-delete the artist-genre relationship when a genre is deleted.
     */
    CONSTRAINT artist_genres_artist_id_fk FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
    CONSTRAINT artist_genres_genre_id_fk FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
);
