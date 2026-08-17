package database

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/KelvinMcClean/palimpsest/hardcover/structs"
)

func (db *DB) SaveBooks(ctx context.Context, books []structs.Book) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	for _, book := range books {
		bookID, err := saveBook(ctx, &tx, book)
		if err != nil {
			log.Fatal("Couldn't insert book")
			return err
		}
		for _, author := range book.Authors {
			authorID, err := saveAuthor(ctx, tx, author)
			if err != nil {
				return err
			}

			err2 := linkBookAuthor(ctx, tx, bookID, authorID, book, author)
			if err2 != nil {
				return err2
			}

			if book.Series.Name != "" {
				err3 := saveAllSeries(ctx, tx, book, bookID)
				if err3 != nil {
					return err3
				}
			}
			err = saveIdentifiers(ctx, book, err, tx, bookID)
			if err2 != nil {
				return err2
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (db *DB) SaveAuthors(ctx context.Context, authors []structs.Author) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, author := range authors {

		authorID, err1 := saveAuthor(ctx, tx, author)
		if err1 != nil {
			return err1
		}

		for _, book := range author.Books {
			bookID, err := saveBook(ctx, &tx, book)
			if err != nil {
				return fmt.Errorf(
					"saving book %d: %w",
					book.ID,
					err,
				)
			}

			err2 := linkBookAuthor(ctx, tx, bookID, authorID, book, author)
			if err2 != nil {
				return err2
			}

			if book.Series.Name != "" {
				err3 := saveAllSeries(ctx, tx, book, bookID)
				if err3 != nil {
					return err3
				}
			}
			err = saveIdentifiers(ctx, book, err, tx, bookID)
			if err2 != nil {
				return err2
			}

		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func saveIdentifiers(ctx context.Context, book structs.Book, err error, tx pgx.Tx, bookID int64) error {
	if book.Editions != nil {
		for _, edition := range book.Editions {
			_, err = tx.Exec(ctx, `
						INSERT INTO hardcover.identifiers (book_id, type, value)
						VALUES ($1, $2, $3)
						ON CONFLICT DO NOTHING
					`, bookID, edition.Type, edition.Value)
			if err != nil {
				return fmt.Errorf(
					"saving edition %s for book %d: %w",
					edition.Type,
					book.ID,
					err,
				)
			}
		}
	}
	_, err = tx.Exec(ctx, `
						INSERT INTO hardcover.identifiers (book_id, type, value)
						VALUES ($1, $2, $3)
						ON CONFLICT DO NOTHING
					`, bookID, "hardcover-id", strconv.Itoa(book.ID),
	)
	return err
}

func saveAllSeries(ctx context.Context, tx pgx.Tx, book structs.Book, bookID int64) error {
	seriesID, err := saveSeries(ctx, &tx, book.Series)
	if err != nil {
		return fmt.Errorf(
			"saving series %d: %w",
			book.Series.ID,
			err,
		)
	}
	err3 := linkBookSeries(ctx, err, tx, bookID, seriesID, book)
	if err3 != nil {
		return err3
	}
	return nil
}

func linkBookSeries(ctx context.Context, err error, tx pgx.Tx, bookID int64, seriesID int64, book structs.Book) error {
	_, err = tx.Exec(ctx, `
					INSERT INTO hardcover.book_series (book_id, series_id, position)
					VALUES ($1, $2, $3)
					ON CONFLICT DO NOTHING
				`, bookID, seriesID, book.Series.Position)

	if err != nil {
		return fmt.Errorf(
			"linking book %d to series %d: %w",
			book.ID,
			seriesID,
			err,
		)
	}
	return nil
}

func linkBookAuthor(ctx context.Context, tx pgx.Tx, bookID int64, authorID int64, book structs.Book, author structs.Author) error {
	if _, err := tx.Exec(ctx, `
				INSERT INTO hardcover.book_authors (book_id, author_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, bookID, authorID); err != nil {
		return fmt.Errorf(
			"linking book %d to author %d: %w",
			book.ID,
			author.ID,
			err,
		)
	}
	return nil
}

func saveSeries(ctx context.Context, tx *pgx.Tx, series structs.BookSeries) (int64, error) {
	var seriesID int64
	err := (*tx).QueryRow(ctx, `
		INSERT INTO hardcover.series (hardcover_id, name)
		VALUES ($1, $2)
		ON CONFLICT (hardcover_id)
		DO UPDATE SET
			hardcover_id = EXCLUDED.hardcover_id,
			name = EXCLUDED.name,
			updated_at = NOW()
		RETURNING id
	`, series.ID, series.Name).Scan(&seriesID)

	if err != nil {
		return 0, err
	}

	return seriesID, nil
}

func saveBook(
	ctx context.Context,
	tx *pgx.Tx,
	book structs.Book,
) (int64, error) {
	var bookID int64

	err := (*tx).QueryRow(ctx, `
		INSERT INTO hardcover.books (hardcover_id, title, slug, subtitle)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (hardcover_id)
		DO UPDATE SET
			hardcover_id = EXCLUDED.hardcover_id,
			title = EXCLUDED.title,
			slug = EXCLUDED.slug,
			subtitle = EXCLUDED.subtitle,
			updated_at = NOW()
		RETURNING id
	`, book.ID, book.Title, book.Slug, book.Subtitle).Scan(&bookID)

	if err != nil {
		return 0, err
	}

	return bookID, nil
}

func saveAuthor(ctx context.Context, tx pgx.Tx, author structs.Author) (int64, error) {
	var authorID int64

	err := tx.QueryRow(ctx, `
			INSERT INTO hardcover.authors (hardcover_id, name)
			VALUES ($1, $2)
			ON CONFLICT (hardcover_id)
			DO UPDATE SET
				name = EXCLUDED.name,
				updated_at = NOW()
			RETURNING id
		`, author.ID, author.Name).Scan(&authorID)

	if err != nil {
		return -1, fmt.Errorf(
			"saving author %d (%s): %w",
			author.ID,
			author.Name,
			err,
		)
	}
	return authorID, nil
}
