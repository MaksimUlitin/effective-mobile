
CREATE TABLE songs (
                       id SERIAL PRIMARY KEY,
                       "group" VARCHAR(255) NOT NULL,
                       song VARCHAR(255) NOT NULL,
                       release_date DATE,
                       text TEXT,
                       link VARCHAR(2083)
);

-- Добавляем индекс для быстрого поиска по группе и названию песни
CREATE INDEX idx_group_song ON songs ("group", song);