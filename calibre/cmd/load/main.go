package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/KelvinMcClean/palimpsest/calibre/calibre"
	"github.com/KelvinMcClean/palimpsest/database/calibreDb"

	_ "modernc.org/sqlite"
)

func main() {

	db, err := sql.Open("sqlite", "file:metadata.db?mode=ro")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	books := getBooks(db)
	formats := getFormats(db)
	authors := getAuthors(db)
	identifiers := getIdentifiers(db)
	series := getSeries(db)
	bookSeries := getBookSeries(db)
	bookAuthors := getBookAuthors(db)

	fmt.Println(books)
	fmt.Println(formats)
	fmt.Println(authors)
	fmt.Println(identifiers)
	fmt.Println(series)
	fmt.Println(bookSeries)
	fmt.Println(bookAuthors)

	ctx := context.Background()

	databaseConn, err := calibreDb.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer databaseConn.Close(ctx)
	err = databaseConn.SaveCalibreData(ctx, books, authors, formats, identifiers, bookAuthors, bookSeries)
	if err != nil {
		panic(err)
	}
}

func getBooks(db *sql.DB) []calibre.Book {
	rows, err := db.Query("SELECT id, title, path FROM books")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var books []calibre.Book

	for rows.Next() {
		var book calibre.Book
		if err := rows.Scan(&book.CalibreID, &book.Title, &book.Path); err != nil {
			panic(err)
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return books
}

func getFormats(db *sql.DB) []calibre.Format {
	rows, err := db.Query("SELECT book, format, name, uncompressed_size FROM data")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var formats []calibre.Format
	for rows.Next() {
		var format calibre.Format
		if err := rows.Scan(&format.BookID, &format.Format, &format.Name, &format.Size); err != nil {
			panic(err)
		}
		formats = append(formats, format)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return formats
}

func getAuthors(db *sql.DB) []calibre.Author {
	rows, err := db.Query("SELECT id, name, sort FROM authors")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var authors []calibre.Author
	for rows.Next() {
		var author calibre.Author
		if err := rows.Scan(&author.ID, &author.Name, &author.Sort); err != nil {
			panic(err)
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return authors
}

func getIdentifiers(db *sql.DB) []calibre.Identifier {
	rows, err := db.Query("SELECT book, type, val FROM identifiers")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var identifiers []calibre.Identifier
	for rows.Next() {
		var identifier calibre.Identifier
		if err := rows.Scan(&identifier.BookID, &identifier.Type, &identifier.Value); err != nil {
			panic(err)
		}
		identifiers = append(identifiers, identifier)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return identifiers
}

func getSeries(db *sql.DB) []calibre.Series {
	rows, err := db.Query("SELECT id, name, sort FROM series")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var series []calibre.Series
	for rows.Next() {
		var s calibre.Series
		if err := rows.Scan(&s.ID, &s.Name, &s.Sort); err != nil {
			panic(err)
		}
		series = append(series, s)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return series
}

func getBookSeries(db *sql.DB) []calibre.BookSeries {
	rows, err := db.Query("SELECT book, series, b.series_index FROM books_series_link bs INNER JOIN books b ON b.id = bs.book")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var bookSeries []calibre.BookSeries
	for rows.Next() {
		var bs calibre.BookSeries
		if err := rows.Scan(&bs.BookID, &bs.SeriesID, &bs.Position); err != nil {
			panic(err)
		}
		bookSeries = append(bookSeries, bs)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return bookSeries
}

func getBookAuthors(db *sql.DB) []calibre.BookAuthor {
	rows, err := db.Query("SELECT book, author FROM books_authors_link")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var bookAuthors []calibre.BookAuthor
	for rows.Next() {
		var ba calibre.BookAuthor
		if err := rows.Scan(&ba.BookID, &ba.AuthorID); err != nil {
			panic(err)
		}
		bookAuthors = append(bookAuthors, ba)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return bookAuthors
}
