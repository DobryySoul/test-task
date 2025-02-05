CREATE TABLE 
    IF NOT EXISTS songs (
        id SERIAL PRIMARY KEY,
        group_name VARCHAR(255) NOT NULL,
        song_name VARCHAR(255) NOT NULL,
        release_date VARCHAR(50),
        text TEXT,
        link VARCHAR(255)
    );