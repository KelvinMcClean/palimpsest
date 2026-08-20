CREATE table hardcover.ingestions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    route TEXT not null,
    created_at timestamp default now(),
    updated_at timestamp default now()
 );

 ALTER TABLE hardcover.books ADD COLUMN ingestion_id BIGINT REFERENCES hardcover.ingestions(id) ON DELETE SET NULL;