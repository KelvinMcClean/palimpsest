
CREATE TABLE calibre.series (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    calibre_id BIGINT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    sort TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE calibre.book_series (
    book_id BIGINT NOT NULL REFERENCES calibre.books(id) ON DELETE CASCADE,
    series_id BIGINT NOT NULL REFERENCES calibre.series(id) ON DELETE CASCADE,
    index NUMERIC(10, 5) NOT NULL,

    PRIMARY KEY (book_id, series_id)
);

CREATE UNIQUE INDEX idx_calibre_book_series_series_id
    ON calibre.book_series(series_id, book_id);

CREATE UNIQUE INDEX idx_series_id
    ON calibre.series(calibre_id);

CREATE INDEX idx_calibre_book_series_book_id
    ON calibre.book_series(book_id);

CREATE INDEX idx_calibre_book_id
    ON calibre.books(calibre_id);
