package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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

func (c *APIClient) GetToRead() []Book {
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

	var books []Book
	for _, meItem := range responseObj.Data.Me {
		for _, userBook := range meItem.UserBooks {
			book := userBook.Book
			var authors []string
			for _, contributor := range book.Contributors {
				authors = append(authors, contributor.Author.Name)
			}
			books = append(books, Book{
				ID:      book.ID,
				Title:   book.Title,
				Authors: authors,
			})
		}
	}
	return books
}

func (c *APIClient) GetFollowedAuthors(config tomlConfig.Config) {
	var query = authorQuery
	query = strings.ReplaceAll(query, "--TOFOLLOWED_AUTHORS_LIST_ID--", strconv.Itoa(config.Hardcover.FollowedAuthorsListID)) // Replace with the actual list ID for followed authors
	resp, err := c.makeRequest(http.MethodPost, strings.NewReader(fmt.Sprintf(`{"query": %q}`, query)))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	var responseObj map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Response object: ", responseObj)
}
