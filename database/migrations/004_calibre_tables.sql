CREATE SCHEMA calibre;

CREATE TABLE calibre.books (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    calibre_id BIGINT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE calibre.identifiers (
    book_id BIGINT NOT NULL REFERENCES calibre.books(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    value TEXT NOT NULL,

    PRIMARY KEY (book_id, type, value)
);

CREATE TABLE calibre.formats(
    book_id BIGINT NOT NULL REFERENCES calibre.books(id) ON DELETE CASCADE,
    format TEXT NOT NULL,
    name TEXT NOT NULL,
    size BIGINT NOT NULL,

    PRIMARY KEY (book_id, format)
);


CREATE INDEX idx_identifiers_book_id
    ON calibre.identifiers(book_id);

CREATE INDEX idx_formats_book_id
    ON calibre.formats(book_id);