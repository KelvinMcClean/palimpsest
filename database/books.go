package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/KelvinMcClean/palimpsest/hardcover/structs"
)

func (db *DB) SaveAuthors(ctx context.Context, authors []structs.Author) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, author := range authors {
		var authorID int64

		err := tx.QueryRow(ctx, `
			INSERT INTO authors (hardcover_id, name)
			VALUES ($1, $2)
			ON CONFLICT (hardcover_id)
			DO UPDATE SET
				name = EXCLUDED.name,
				updated_at = NOW()
			RETURNING id
		`, author.ID, author.Name).Scan(&authorID)

		if err != nil {
			return fmt.Errorf(
				"saving author %d (%s): %w",
				author.ID,
				author.Name,
				err,
			)
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

			if _, err := tx.Exec(ctx, `
				INSERT INTO book_authors (book_id, author_id)
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
			if book.Series.Name != "" {
				seriesID, err := saveSeries(ctx, &tx, book.Series)
				if err != nil {
					return fmt.Errorf(
						"saving series %d: %w",
						book.Series.ID,
						err,
					)
				}
				_, err = tx.Exec(ctx, `
					INSERT INTO book_series (book_id, series_id, position)
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
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func saveSeries(ctx context.Context, tx *pgx.Tx, series structs.BookSeries) (int64, error) {
	var seriesID int64
	err := (*tx).QueryRow(ctx, `
		INSERT INTO series (hardcover_id, name)
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
		INSERT INTO books (hardcover_id, title, slug, subtitle)
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
