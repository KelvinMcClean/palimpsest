package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/KelvinMcClean/palimpsest/hardcover/structs"
	"github.com/KelvinMcClean/palimpsest/hardcover/tomlConfig"
)

type APIClient struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	token      string
}

func NewAPIClient(baseURL string, config tomlConfig.Config) (*APIClient, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &APIClient{
		BaseURL: parsedURL,
		token:   config.Hardcover.Key,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *APIClient) makeRequest(method string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.BaseURL.String(), body)
	if err != nil {
		log.Fatal("error creating request", err)
		return nil, err
	}

	req.Header.Set("Authorization", c.token) // Token includes "Bearer " by default on copy-paste from the hardcover API generation point.
	req.Header.Set("Content-Type", "application/json")

	log.Println(req.RequestURI, req.Header, req.Body)

	return c.HTTPClient.Do(req)
}

func (c *APIClient) GetToRead(config tomlConfig.Config) ([]structs.Book, bool) {
	query := toReadQuery

	limit := 100
	offset := 0

	newBooks := true
	var books []structs.Book

	for newBooks {
		resp, err := c.makeRequest(http.MethodPost, strings.NewReader(fmt.Sprintf(`{"query": %q, "variables": {"limit": %d, "offset": %d}}`,
			query, limit, offset)))
		if err != nil {
			log.Fatal(err)
			return nil, false
		}
		log.Println("Response Status:", resp.Status)
		defer resp.Body.Close()
		var responseObj wantToReadResponse
		err = json.NewDecoder(resp.Body).Decode(&responseObj)
		if err != nil {
			log.Fatal(err)
			return nil, false
		}
		if len(responseObj.Data.Me) == 0 {
			return books, true
		}
		log.Println("Response object: ", responseObj)
		books, newBooks = parseBooksFromWantToRead(responseObj, books)
		offset += limit
	}
	return books, true

}

func parseBooksFromWantToRead(responseObj wantToReadResponse, books []structs.Book) ([]structs.Book, bool) {
	for _, meItem := range responseObj.Data.Me {
		if len(meItem.UserBooks) == 0 {
			return books, false
		}
		for _, userBook := range meItem.UserBooks {
			book := userBook.Book
			var authors []structs.Author
			for _, contributor := range book.Contributors {
				authors = append(authors, structs.Author{
					Name: contributor.Author.Name,
					ID:	  contributor.Author.ID})
			}
			var bookObj = structs.Book{
					ID:       book.ID,
					Title:    book.Title,
					Slug:     book.Slug,
					Authors:  authors,
					Subtitle: book.Subtitle,
				}
				// Get the series name from the series matching the FeaturedBookSeriesID
				for _, series := range book.BookSeries {
					if series.ID == book.FeaturedBookSeriesID {
						bookObj.Series.Name = series.Series.Name
						bookObj.Series.Position = series.Position
						bookObj.Series.ID = series.Series.ID
						break
					}
				}
				for _, edition := range book.Editions {

					if edition.ISBN10 != "" {
						bookObj.Editions = append(bookObj.Editions, structs.Edition{Type: "ISBN10", Value: edition.ISBN10})

					}
					if edition.ISBN13 != "" {
						bookObj.Editions = append(bookObj.Editions, structs.Edition{Type: "ISBN13", Value: edition.ISBN13})
					}
					bookObj.Editions = append(bookObj.Editions, structs.Edition{Type: "hardcover-id", Value: strconv.Itoa(edition.ID)})
				}
				books = append(books, bookObj)
		}
	}
	return books, true
}

func (c *APIClient) GetFollowedAuthors(config tomlConfig.Config) ([]structs.Author, bool) {
	limit := 100
	offset := 0
	minUsersCount := config.Hardcover.MinUsersCount
	// We loop on the Books on the Author object, so we don't need to paginate the authors themselves, just the books for each author.
	var totalBooks = math.MaxInt64
	var authors []structs.Author
	for totalBooks >= limit {
		var moreAuthors, moreOk = c.getFollowedAuthors(config, limit, offset, minUsersCount)
		if !moreOk {
			log.Fatal("Error getting followed authors")
			return nil, false
		}
		offset += limit

		authors = structs.MergeAuthors(authors, moreAuthors)
		totalBooks = getTotalBooks(moreAuthors)
	}

	return authors, true
}

func getTotalBooks(authors []structs.Author) int {
	var totalBooks = 0
	for _, author := range authors {
		totalBooks += len(author.Books)
	}
	return totalBooks
}

func (c *APIClient) getFollowedAuthors(config tomlConfig.Config, limit int, offset int, minUsersCount int) ([]structs.Author, bool) {
	var query = authorQuery
	resp, err := c.makeRequest(http.MethodPost, strings.NewReader(fmt.Sprintf(`{"query": %q, "variables": {"id": %d, "limit": %d, "offset": %d, "minUsersCount": %d}}`, query,
		config.Hardcover.FollowedAuthorsListID, limit, offset, minUsersCount)))
	if err != nil {
		log.Fatal(err)
		return nil, false
	}
	defer resp.Body.Close()
	var responseObj authorResponse
	err = json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		log.Fatal(err)
		return nil, false
	}
	log.Println("Response object: ", responseObj)
	var authors []structs.Author
	for _, filterList := range responseObj.Data.FilterLists {
		for _, entity := range filterList.FilterListEntities {
			var authorObj = structs.Author{
				ID:   entity.Author.ID,
				Name: entity.Author.Name,
			}
			for _, contribution := range entity.Author.Contributions {
				var bookObj = structs.Book{
					ID:       contribution.Book.ID,
					Title:    contribution.Book.Title,
					Slug:     contribution.Book.Slug,
					Authors:  []structs.Author{authorObj},
					Subtitle: contribution.Book.Subtitle,
				}
				// Get the series name from the series matching the FeaturedBookSeriesID
				for _, series := range contribution.Book.BookSeries {
					if series.ID == contribution.Book.FeaturedBookSeriesID {
						bookObj.Series.Name = series.Series.Name
						bookObj.Series.Position = series.Position
						bookObj.Series.ID = series.Series.ID
						break
					}
				}
				for _, edition := range contribution.Book.Editions {

					if edition.ISBN10 != "" {
						bookObj.Editions = append(bookObj.Editions, structs.Edition{Type: "ISBN10", Value: edition.ISBN10})

					}
					if edition.ISBN13 != "" {
						bookObj.Editions = append(bookObj.Editions, structs.Edition{Type: "ISBN13", Value: edition.ISBN13})
					}
					bookObj.Editions = append(bookObj.Editions, structs.Edition{Type: "hardcover-id", Value: strconv.Itoa(edition.ID)})
				}

				authorObj.Books = append(authorObj.Books, bookObj)
			}
			authors = append(authors, authorObj)
		}
	}

	return authors, true
}
