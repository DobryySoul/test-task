CREATE TABLE Songs (
    song_id SERIAL PRIMARY KEY,
    song_name VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    song_text TEXT,
    link VARCHAR(255),
    artist_id SERIAL NOT NULL,
    FOREIGN KEY (artist_id) REFERENCES Artists(artist_id)
);
