CREATE SCHEMA hardcover;
CREATE SCHEMA app;

ALTER TABLE public.authors
    SET SCHEMA hardcover;

ALTER TABLE public.book_authors
    SET SCHEMA hardcover;

ALTER TABLE public.book_series
    SET SCHEMA hardcover;

ALTER TABLE public.books
    SET SCHEMA hardcover;

ALTER TABLE public.series
    SET SCHEMA hardcover;