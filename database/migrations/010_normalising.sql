CREATE TABLE calibre.search_books (
    id BIGINT REFERENCES calibre.books(id) ON DELETE CASCADE PRIMARY KEY,
    title TEXT NOT NULL,
    normalized_title TEXT NOT NULL,
    normalized_title_base TEXT NOT NULL,
    subtitle TEXT,
    normalized_subtitle TEXT,
    normalized_subtitle_base TEXT,
    authors TEXT[] NOT NULL,
    normalized_authors TEXT[] NOT NULL,
    series TEXT[] NOT NULL,
    normalized_series TEXT[] NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE hardcover.search_books (
    id BIGINT REFERENCES hardcover.books(id) ON DELETE CASCADE PRIMARY KEY,
    title TEXT NOT NULL,
    normalized_title TEXT NOT NULL,
    normalized_title_base TEXT NOT NULL,
    subtitle TEXT,
    normalized_subtitle TEXT,
    normalized_subtitle_base TEXT,
    authors TEXT[] NOT NULL,
    normalized_authors TEXT[] NOT NULL,
    series TEXT[] NOT NULL,
    normalized_series TEXT[] NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);