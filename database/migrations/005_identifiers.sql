CREATE TABLE hardcover.identifiers (
    book_id BIGINT NOT NULL REFERENCES hardcover.books(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    value TEXT NOT NULL,

    PRIMARY KEY (book_id, type, value)
);

CREATE INDEX idx_identifiers_book_id
    ON hardcover.identifiers(book_id);