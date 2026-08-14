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
	Series  struct {
		Name     string  `json:"name"`
		Position float64 `json:"position"`
	} `json:"series"`
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
	if b.Series.Name != "" {
		formatPosition := strconv.FormatFloat(b.Series.Position, 'f', -1, 64)
		return fmt.Sprintf("  Book: %s (ID: %d) - Series: %s (Position: %s)\n", b.Title, b.ID, b.Series.Name, formatPosition)
	} else {
		return fmt.Sprintf("  Book: %s (ID: %d)\n", b.Title, b.ID)
	}
}

func (a Author) String() string {
	var result string
	result += fmt.Sprintf("Author ID: %d, Name: %s\n", a.ID, a.Name)
		for _, book := range a.Books {
			result += book.String()
		}
	return result
}