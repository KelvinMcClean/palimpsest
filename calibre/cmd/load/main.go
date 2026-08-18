package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
	"calibre"
)

func main() {

	db, err := sql.Open("sqlite", "file:metadata.db?mode=ro")
	if err != nil {
		panic(err)
	}

	books := getBooks(db)
	formats := getFormats(db)
	authors := getAuthors(db)
	identifiers := getIdentifiers(db)
	fmt.Println(books)
	fmt.Println(formats)
	fmt.Println(authors)
	fmt.Println(identifiers)

	defer db.Close()

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