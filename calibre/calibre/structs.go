package calibre

import (
	"github.com/KelvinMcClean/palimpsest/normalize"
)

type Book struct {
	CalibreID int
	Title     string
	Path      string
}

type Format struct {
	BookID int
	Format string
	Name   string
	Size   int64
}

type Author struct {
	ID   int
	Name string
	Sort string
}

type Identifier struct {
	BookID int
	Type   string
	Value  string
}

type Series struct {
	ID       int
	Name     string
	Sort     string
}

type BookAuthor struct {
	BookID   int
	AuthorID int
}

type BookSeries struct {
	BookID   int
	SeriesID int
	Position int
}


type NormalizedBook struct {
	ID                     int
	Title                  string
	NormalizedTitle        string
	NormalizedTitleBase    string
	Authors                []string
	NormalizedAuthors      []string
	Series                 []string
	NormalizedSeries       []string
}

func JoinNormalizedBooks(normalizedBooks []NormalizedBook) []NormalizedBook {
	joined := map[int]NormalizedBook{}
	for _, book := range normalizedBooks {
		if existing, ok := joined[book.ID]; ok {
			existing.Authors = append(existing.Authors, book.Authors...)
			existing.NormalizedAuthors = append(existing.NormalizedAuthors, book.NormalizedAuthors...)
			existing.Series = append(existing.Series, book.Series...)
			existing.NormalizedSeries = append(existing.NormalizedSeries, book.NormalizedSeries...)
			joined[book.ID] = existing
		} else {
			joined[book.ID] = book
		}
	}

	// Convert the map back to a slice
	result := []NormalizedBook{}
	for _, book := range joined {
		result = append(result, book)
	}
	return result
}

func (b NormalizedBook) Normalize() NormalizedBook {
	var normalizedBook NormalizedBook
	normalizedBook.ID = b.ID
	normalizedBook.Title = b.Title
	normalizedBook.NormalizedTitle = normalize.Normalize(b.Title)
	normalizedBook.NormalizedTitleBase = normalize.NormalizeBase(b.Title)
	normalizedBook.Authors = b.Authors
	for _, author := range b.Authors {
		normalizedBook.NormalizedAuthors = append(normalizedBook.NormalizedAuthors, normalize.NormalizeAuthorName(author))
	}
	if (len(b.Series) == 0) {
		b.Series = []string{}
		b.Series = append(b.Series, "")
	}
	normalizedBook.Series = b.Series
	
	for _, series := range b.Series {
		normalizedBook.NormalizedSeries = append(normalizedBook.NormalizedSeries, normalize.Normalize(series))
	}
	return normalizedBook
}


func (b Book) Normalize(series []Series, authors []Author) NormalizedBook {
	var normalizedBook NormalizedBook
	normalizedBook.ID = b.CalibreID
	normalizedBook.Title = b.Title
	normalizedBook.NormalizedTitle = normalize.Normalize(b.Title)
	normalizedBook.NormalizedTitleBase = normalize.NormalizeBase(b.Title)
	var authorsS []string
	var normalizedAuthors []string
	for _, author := range authors {
		authorsS = append(authorsS, author.Name)
		normalizedAuthors = append(normalizedAuthors, normalize.NormalizeAuthorName(author.Name))
	}
	normalizedBook.Authors = authorsS
	normalizedBook.NormalizedAuthors = normalizedAuthors
	var seriesS []string
	var normalizedSeries []string
	for _, s := range series {
		seriesS = append(seriesS, s.Name)
		normalizedSeries = append(normalizedSeries, normalize.Normalize(s.Name))
	}
	normalizedBook.Series = seriesS
	normalizedBook.NormalizedSeries = normalizedSeries
	return normalizedBook
}
