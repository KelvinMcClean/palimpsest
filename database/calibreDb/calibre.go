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

func (cdb *cDB) saveCalibreBooks(tx pgx.Tx, ctx context.Context, books []calibre.Book) error {
	for _, book := range books {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO calibre.books (calibre_id, title, path)
	 			VALUES ($1, $2, $3)
				ON CONFLICT DO NOTHING`,
			book.CalibreID,
			book.Title,
			book.Path,
		)
		if err != nil {
			return fmt.Errorf("inserting books: %w", err)
		}
	}
	return nil
}

func (cdb *cDB) saveCalibreAuthors(tx pgx.Tx, ctx context.Context, authors []calibre.Author) error {
	

	for _, author := range authors {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO calibre.authors (calibre_id, name, sort)
	 			VALUES ($1, $2, $3)
				ON CONFLICT DO NOTHING`,
			author.ID,
			author.Name,
			author.Sort,
		)
		if err != nil {
			return fmt.Errorf("inserting authors: %w", err)
		}
	}

	return nil
}

func (cdb *cDB) saveCalibreFormats(tx pgx.Tx, ctx context.Context, formats []calibre.Format) error {
	
	for _, format := range formats {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO calibre.formats (book_id, format, name, size)
     			SELECT id, $2, $3, $4
     			FROM calibre.books
     			WHERE calibre_id = $1
				ON CONFLICT DO NOTHING`,
			format.BookID,
			format.Format,
			format.Name,
			format.Size,
		)
		if err != nil {
			return fmt.Errorf("inserting formats: %w", err)
		}
	}
	return nil
}

func (cdb *cDB) saveCalibreIdentifiers(tx pgx.Tx, ctx context.Context, identifiers []calibre.Identifier) error {
	for _, identifier := range identifiers {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO calibre.identifiers (book_id, type, value)
     			SELECT id, $2, $3
     			FROM calibre.books
     			WHERE calibre_id = $1
				ON CONFLICT DO NOTHING`,
			identifier.BookID,
			identifier.Type,
			identifier.Value,
		)
		if err != nil {
			return fmt.Errorf("inserting identifiers: %w", err)
		}
	}
	return nil
}

func (cdb *cDB) saveCalibreBookAuthors(tx pgx.Tx, ctx context.Context, bookAuthors []calibre.BookAuthor) error {
	// tx, err := cdb.db.Conn.Begin(ctx)
	// if err != nil {
	// 	return fmt.Errorf("starting transaction: %w", err)
	// }
	// defer tx.Rollback(ctx)
	for _, bookAuthor := range bookAuthors {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO calibre.book_authors (book_id, author_id)
     			SELECT b.id, a.id
     			FROM calibre.books b
     			JOIN calibre.authors a ON a.calibre_id = $2
     			WHERE b.calibre_id = $1
				ON CONFLICT DO NOTHING`,
			bookAuthor.BookID,
			bookAuthor.AuthorID,
		)
		if err != nil {
			return fmt.Errorf("inserting book authors: %w", err)
		}
	}
	// if err := tx.Commit(ctx); err != nil {
	// 	return fmt.Errorf("committing transaction: %w", err)
	// }
	return nil
}


func (cdb *cDB) saveCalibreSeries(tx pgx.Tx, ctx context.Context, series []calibre.Series) error {
	for _, s := range series {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO calibre.series (calibre_id, name, sort)
	 			VALUES ($1, $2, $3)
				ON CONFLICT DO NOTHING`,
				s.ID,
				s.Name,
				s.Sort,
		)
		if err != nil {
			return fmt.Errorf("inserting series: %w", err)
		}
	}
	// if err := tx.Commit(ctx); err != nil {
	// 	return fmt.Errorf("committing transaction: %w", err)
	// }
	return nil
}


func (cdb *cDB) saveCalibreBookSeries(tx pgx.Tx, ctx context.Context, bookSeries []calibre.BookSeries) error {
	for _, bs := range bookSeries {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO calibre.book_series (book_id, series_id, index)
	 			SELECT b.id, s.id, $3
	 			FROM calibre.books b
	 			JOIN calibre.series s ON s.calibre_id = $2
	 			WHERE b.calibre_id = $1
				ON CONFLICT DO NOTHING`,
			bs.BookID,
			bs.SeriesID,
			bs.Position,
		)
		if err != nil {
			return fmt.Errorf("inserting book series: %w", err)
		}
	}
	// if err := tx.Commit(ctx); err != nil {
	// 	return fmt.Errorf("committing transaction: %w", err)
	// }
	return nil
}

func (cdb *cDB) TruncateCalibreTables(ctx context.Context) error {
	tx, err := cdb.db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	tables := []string{
		"calibre.book_series",
		"calibre.book_authors",
		"calibre.identifiers",
		"calibre.formats",
		"calibre.books",
		"calibre.authors",
		"calibre.series",
	}
	for _, table := range tables {
		_, err := tx.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("truncating table %s: %w", table, err)
		}
	}
	return tx.Commit(ctx)
}

func (cdb *cDB) SaveCalibreData(ctx context.Context, books []calibre.Book, authors []calibre.Author, formats []calibre.Format, identifiers []calibre.Identifier, bookAuthors []calibre.BookAuthor, bookSeries []calibre.BookSeries, series []calibre.Series	) error {

	tx, err := cdb.db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := cdb.saveCalibreBooks(tx, ctx, books); err != nil {
		return fmt.Errorf("saving books: %w", err)
	}
	if err := cdb.saveCalibreAuthors(tx, ctx, authors); err != nil {
		return fmt.Errorf("saving authors: %w", err)
	}
	if err := cdb.saveCalibreFormats(tx, ctx, formats); err != nil {
		return fmt.Errorf("saving formats: %w", err)
	}
	if err := cdb.saveCalibreIdentifiers(tx, ctx, identifiers); err != nil {
		return fmt.Errorf("saving identifiers: %w", err)
	}
	if err := cdb.saveCalibreBookAuthors(tx, ctx, bookAuthors); err != nil {
		return fmt.Errorf("saving book authors: %w", err)
	}
	if err := cdb.saveCalibreSeries(tx, ctx, series); err != nil {
		return fmt.Errorf("saving series: %w", err)
	}
	if err := cdb.saveCalibreBookSeries(tx, ctx, bookSeries); err != nil {
		return fmt.Errorf("saving book series: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
