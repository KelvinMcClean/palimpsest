package structs

import (
	"fmt"
	"strconv"
)
func MergeAuthors(a []Author, b []Author) []Author {
	authorMap := make(map[int]Author)
	for _, author := range a {
		authorMap[author.ID] = author
	}
	for _, author := range b {
		if existingAuthor, ok := authorMap[author.ID]; ok {
			existingAuthor.Books = append(existingAuthor.Books, author.Books...)
			authorMap[author.ID] = existingAuthor
		} else {
			authorMap[author.ID] = author
		}
	}
	mergedAuthors := make([]Author, 0, len(authorMap))
	for _, author := range authorMap {
		mergedAuthors = append(mergedAuthors, author)
	}
	return mergedAuthors
}

type Author struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Books []Book `json:"books"`
}

type Book struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
	Slug    string   `json:"slug"`
	Subtitle string   `json:"subtitle"`
	Series  BookSeries `json:"series"`
	Editions []Edition `json:"editions"`
}

type Edition struct {
	Type string
	Value string
}

type BookSeries struct {
	ID       int     `json:"id"`
	Position float64 `json:"position"`
	Name     string  `json:"name"`
}

func AuthorsToBooks(authors []Author) []Book {
	var books []Book
	for _, author := range authors {
		for _, book := range author.Books {
			book.Authors = append(book.Authors, author.Name)
			books = append(books, book)
		}
	}
	return books
}

func (b Book) String() string {
	seriesInfo := ""
	editionInfo := ""
	if b.Series.Name != "" {
		formatPosition := strconv.FormatFloat(b.Series.Position, 'f', -1, 64)
		seriesInfo = fmt.Sprintf(" - Series: %s (Position: %s) (ID: %d)", b.Series.Name, formatPosition, b.Series.ID)
	}
	if (len(b.Editions) > 0) {
		editionInfo = " - Editions: "
		for _, edition := range b.Editions {
			editionInfo += fmt.Sprintf("%s: %s; ", edition.Type, edition.Value)
		}
	}
	return fmt.Sprintf("  Book: %s (ID: %d)%s%s\n        %s\n", b.Title, b.ID, seriesInfo, editionInfo, b.Slug)
}

func (a Author) String() string {
	var result string
	result += fmt.Sprintf("Author ID: %d, Name: %s\n", a.ID, a.Name)
		for _, book := range a.Books {
			result += book.String()
		}
	return result
}