CREATE TABLE songs (
                       id SERIAL PRIMARY KEY,
                       band TEXT NOT NULL,
                       title TEXT NOT NULL,
                       release_date TIMESTAMP NOT NULL,
                       text TEXT NOT NULL,
                       link TEXT NOT NULL
);