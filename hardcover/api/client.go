package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/KelvinMcClean/palimpsest/hardcover/tomlConfig"
	"github.com/KelvinMcClean/palimpsest/hardcover/structs"
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

func (c *APIClient) GetToRead() []structs.Book {
	query := toReadQuery

	resp, err := c.makeRequest(http.MethodPost, strings.NewReader(fmt.Sprintf(`{"query": %q}`, query)))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("Response Status:", resp.Status)
	defer resp.Body.Close()
	var responseObj wantToReadResponse
	err = json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if len(responseObj.Data.Me) == 0 {
		log.Fatal("No response from api")
		return nil
	}
	log.Println("Response object: ", responseObj)

	var books []structs.Book
	for _, meItem := range responseObj.Data.Me {
		for _, userBook := range meItem.UserBooks {
			book := userBook.Book
			var authors []string
			for _, contributor := range book.Contributors {
				authors = append(authors, contributor.Author.Name)
			}
			books = append(books, structs.Book{
				ID:      book.ID,
				Title:   book.Title,
				Authors: authors,
			})
		}
	}
	return books
}

func (c *APIClient) GetFollowedAuthors(config tomlConfig.Config) ([]structs.Author, bool) {
	limit := 10
	offset := 0
	minUsersCount := config.Hardcover.MinUsersCount
	// We loop on the Books on the Author object, so we don't need to paginate the authors themselves, just the books for each author.
	var totalBooks = 999999
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
					ID:    contribution.Book.ID,
					Title: contribution.Book.Title,
				}
				// Get the series name from the series matching the FeaturedBookSeriesID
				for _, series := range contribution.Book.BookSeries {
					if series.ID == contribution.Book.FeaturedBookSeriesID {
						bookObj.Series.Name = series.Series.Name
						bookObj.Series.Position = series.Position
						break
					}
				}
				authorObj.Books = append(authorObj.Books, bookObj)
			}
			authors = append(authors, authorObj)
		}
	}

	return authors, true
}
