package hardcoverDb

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/KelvinMcClean/palimpsest/database"

	"github.com/KelvinMcClean/palimpsest/hardcover/structs"
)
	
type hDB struct {
	db *database.DB
}


func Connect(ctx context.Context) (*hDB, error) {
	hdb := &hDB{}
	var err error
	hdb.db, err = database.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	return hdb, nil
}

func (hdb *hDB) Close(ctx context.Context) error {
	return hdb.db.Close(ctx)
}

func (hdb *hDB) SaveHardcoverBooks(ctx context.Context, books []structs.Book, route string) error {
	tx, err := hdb.db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var ingestionID int64
	err = tx.QueryRow(ctx, "INSERT INTO hardcover.ingestions (route) VALUES ($1) RETURNING id", route).Scan(&ingestionID)
	if err != nil {
		return fmt.Errorf("inserting ingestion route: %w", err)
	}
	for _, book := range books {
		bookID, err := saveHardcoverBook(ctx, &tx, book, ingestionID)
		if err != nil {
			log.Fatal("Couldn't insert book")
			return err
		}
		_, err = saveHardcoverBookNormalized(ctx, &tx, book.Normalize(), bookID, ingestionID)
		if err != nil {
			return fmt.Errorf(
				"saving normalized book %d: %w",
				book.ID,
				err,
			)
		}
		for _, author := range book.Authors {
			authorID, err := saveHardcoverAuthor(ctx, tx, author)
			if err != nil {
				return err
			}

			err2 := linkBookAuthor(ctx, tx, bookID, authorID, book, author)
			if err2 != nil {
				return err2
			}
		}

		if book.Series.Name != "" {
			err3 := saveAllHardcoverSeries(ctx, tx, book, bookID)
			if err3 != nil {
				return err3
			}
		}
		err = saveHardcoverIdentifiers(ctx, book, tx, bookID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (hdb *hDB) SaveHardcoverAuthors(ctx context.Context, authors []structs.Author, route string) error {
	tx, err := hdb.db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)


	var ingestionID int64
	err = tx.QueryRow(ctx, "INSERT INTO hardcover.ingestions (route) VALUES ($1) RETURNING id", route).Scan(&ingestionID)
	if err != nil {
		return fmt.Errorf("inserting ingestion route: %w", err)
	}

	for _, author := range authors {

		authorID, err1 := saveHardcoverAuthor(ctx, tx, author)
		if err1 != nil {
			return err1
		}

		for _, book := range author.Books {
			bookID, err := saveHardcoverBook(ctx, &tx, book, ingestionID)
			if err != nil {
				return fmt.Errorf(
					"saving book %d: %w",
					book.ID,
					err,
				)
			}
			_, err = saveHardcoverBookNormalized(ctx, &tx, book.Normalize(), bookID, ingestionID)
			if err != nil {
				return fmt.Errorf(
					"saving normalized book %d: %w",
					book.ID,
					err,
				)
			}

			err2 := linkBookAuthor(ctx, tx, bookID, authorID, book, author)
			if err2 != nil {
				return err2
			}

			if book.Series.Name != "" {
				err3 := saveAllHardcoverSeries(ctx, tx, book, bookID)
				if err3 != nil {
					return err3
				}
			}
			err = saveHardcoverIdentifiers(ctx, book, tx, bookID)
			if err != nil {
				return err
			}

		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func saveHardcoverIdentifiers(ctx context.Context, book structs.Book,  tx pgx.Tx, bookID int64) error {
	if book.Editions != nil {
		for _, edition := range book.Editions {
			_, err := tx.Exec(ctx, `
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
	_, err := tx.Exec(ctx, `
						INSERT INTO hardcover.identifiers (book_id, type, value)
						VALUES ($1, $2, $3)
						ON CONFLICT DO NOTHING
					`, bookID, "hardcover-id", strconv.Itoa(book.ID),
	)
	return err
}

func saveAllHardcoverSeries(ctx context.Context, tx pgx.Tx, book structs.Book, bookID int64) error {
	seriesID, err := saveHardcoverSeries(ctx, &tx, book.Series)
	if err != nil {
		return fmt.Errorf(
			"saving series %d: %w",
			book.Series.ID,
			err,
		)
	}
	err3 := linkBookSeries(ctx, tx, bookID, seriesID, book)
	if err3 != nil {
		return err3
	}
	return nil
}

func linkBookSeries(ctx context.Context, tx pgx.Tx, bookID int64, seriesID int64, book structs.Book) error {
	_, err := tx.Exec(ctx, `
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

func saveHardcoverSeries(ctx context.Context, tx *pgx.Tx, series structs.BookSeries) (int64, error) {
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

func saveHardcoverBook(
	ctx context.Context,
	tx *pgx.Tx,
	book structs.Book,
	ingestionID int64,
) (int64, error) {
	var bookID int64

	err := (*tx).QueryRow(ctx, `
		INSERT INTO hardcover.books (hardcover_id, title, slug, subtitle, ingestion_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (hardcover_id)
		DO UPDATE SET
			hardcover_id = EXCLUDED.hardcover_id,
			title = EXCLUDED.title,
			slug = EXCLUDED.slug,
			subtitle = EXCLUDED.subtitle,
			ingestion_id = EXCLUDED.ingestion_id,
			updated_at = NOW()
		RETURNING id
	`, book.ID, book.Title, book.Slug, book.Subtitle, ingestionID).Scan(&bookID)

	if err != nil {
		return 0, err
	}

	return bookID, nil
}

func saveHardcoverBookNormalized(
	ctx context.Context,
	tx *pgx.Tx,
	book structs.NormalizedBook,
	id int64,
	ingestionID int64,
) (int64, error) {

	_, err := (*tx).Exec(ctx, `
		INSERT INTO hardcover.search_books (id, 
			title, normalized_title, normalized_title_base,
			subtitle, normalized_subtitle, normalized_subtitle_base, 
			authors, normalized_authors, series, normalized_series)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id)
		DO UPDATE SET
			title = EXCLUDED.title,
			normalized_title = EXCLUDED.normalized_title,
			normalized_title_base = EXCLUDED.normalized_title_base,
			subtitle = EXCLUDED.subtitle,
			normalized_subtitle = EXCLUDED.normalized_subtitle,
			normalized_subtitle_base = EXCLUDED.normalized_subtitle_base,
			authors = EXCLUDED.authors,
			normalized_authors = EXCLUDED.normalized_authors,
			series = EXCLUDED.series,
			normalized_series = EXCLUDED.normalized_series	
	`, 
		id, book.Title, book.NormalizedTitle, book.NormalizedTitleBase,
		book.Subtitle, book.NormalizedSubtitle, book.NormalizedSubtitleBase,
		book.Authors, book.NormalizedAuthors,
		book.Series, book.NormalizedSeries)

	if err != nil {
		return 0, fmt.Errorf(
			"saving normalized book %d: %w",
			id,
			err,
		)
	}

	return id, nil
}

func saveHardcoverAuthor(ctx context.Context, tx pgx.Tx, author structs.Author) (int64, error) {
	var authorID int64
	log.Println("Saving author: ", author.ID, author.Name)
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
