package calibreDb

import (
	"context"
	"fmt"

	"github.com/KelvinMcClean/palimpsest/calibre/calibre"
	"github.com/KelvinMcClean/palimpsest/database"
	"github.com/jackc/pgx/v5"
)

type cDB struct {
	db *database.DB
}

func Connect(ctx context.Context) (*cDB, error) {
	cdb := &cDB{}
	var err error
	cdb.db, err = database.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	return cdb, nil
}

func (cdb *cDB) Close(ctx context.Context) error {
	return cdb.db.Close(ctx)
}

func (cdb *cDB) saveCalibreBooks(ctx context.Context, books []calibre.Book) error {
	tx, err := cdb.db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	rows := make([][]any, 0, len(books))

	for _, book := range books {
		rows = append(rows, []any{
			book.CalibreID,
			book.Title,
			book.Path,
		})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"calibre.books"},
		[]string{"calibre_id", "title", "path"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("inserting books: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (cdb *cDB) saveCalibreAuthors(ctx context.Context, authors []calibre.Author) error {
	tx, err := cdb.db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	rows := make([][]any, 0, len(authors))

	for _, author := range authors {
		rows = append(rows, []any{
			author.ID,
			author.Name,
			author.Sort,
		})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"calibre.authors"},
		[]string{"calibre_id", "name", "sort"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("inserting authors: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (cdb *cDB) saveCalibreFormats(ctx context.Context, formats []calibre.Format) error {
	tx, err := cdb.db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	for _, format := range formats {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO calibre.formats (book_id, format, name, size)
     			SELECT id, $2, $3, $4
     			FROM calibre.books
     			WHERE calibre_id = $1`,
			format.BookID,
			format.Format,
			format.Name,
			format.Size,
		)
		if err != nil {
			return fmt.Errorf("inserting formats: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (cdb *cDB) saveCalibreIdentifiers(ctx context.Context, identifiers []calibre.Identifier) error {
	tx, err := cdb.db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	for _, identifier := range identifiers {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO calibre.identifiers (book_id, type, value)
     			SELECT id, $2, $3
     			FROM calibre.books
     			WHERE calibre_id = $1`,
			identifier.BookID,
			identifier.Type,
			identifier.Value,
		)
		if err != nil {
			return fmt.Errorf("inserting identifiers: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (cdb *cDB) saveCalibreBookAuthors(ctx context.Context, bookAuthors []calibre.BookAuthor) error {
	tx, err := cdb.db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	for _, bookAuthor := range bookAuthors {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO calibre.book_authors (book_id, author_id)
     			SELECT b.id, a.id
     			FROM calibre.books b
     			JOIN calibre.authors a ON a.calibre_id = $2
     			WHERE b.calibre_id = $1`,
			bookAuthor.BookID,
			bookAuthor.AuthorID,
		)
		if err != nil {
			return fmt.Errorf("inserting book authors: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (cdb *cDB) saveCalibreBookSeries(ctx context.Context, bookSeries []calibre.BookSeries) error {
	tx, err := cdb.db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	for _, bs := range bookSeries {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO calibre.book_series (book_id, series_id, position)
	 			SELECT b.id, s.id, $3
	 			FROM calibre.books b
	 			JOIN calibre.series s ON s.calibre_id = $2
	 			WHERE b.calibre_id = $1`,
			bs.BookID,
			bs.SeriesID,
			bs.Position,
		)
		if err != nil {
			return fmt.Errorf("inserting book series: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (cdb *cDB) SaveCalibreData(ctx context.Context, books []calibre.Book, authors []calibre.Author, formats []calibre.Format, identifiers []calibre.Identifier, bookAuthors []calibre.BookAuthor, bookSeries []calibre.BookSeries) error {

	if err := cdb.saveCalibreBooks(ctx, books); err != nil {
		return fmt.Errorf("saving books: %w", err)
	}
	if err := cdb.saveCalibreAuthors(ctx, authors); err != nil {
		return fmt.Errorf("saving authors: %w", err)
	}
	if err := cdb.saveCalibreFormats(ctx, formats); err != nil {
		return fmt.Errorf("saving formats: %w", err)
	}
	if err := cdb.saveCalibreIdentifiers(ctx, identifiers); err != nil {
		return fmt.Errorf("saving identifiers: %w", err)
	}
	if err := cdb.saveCalibreBookAuthors(ctx, bookAuthors); err != nil {
		return fmt.Errorf("saving book authors: %w", err)
	}
	if err := cdb.saveCalibreBookSeries(ctx, bookSeries); err != nil {
		return fmt.Errorf("saving book series: %w", err)
	}
	return nil
}
