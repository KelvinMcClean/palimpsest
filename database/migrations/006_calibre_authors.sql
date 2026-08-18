
CREATE TABLE calibre.authors (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    calibre_id BIGINT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    sort TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

create table calibre.book_authors (
    book_id BIGINT NOT NULL REFERENCES calibre.books(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES calibre.authors(id) ON DELETE CASCADE,

    PRIMARY KEY (book_id, author_id)
);

CREATE INDEX idx_calibre_book_authors_author_id
    ON calibre.book_authors(author_id);