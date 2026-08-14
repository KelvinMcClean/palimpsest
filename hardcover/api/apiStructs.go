package api

type contributor struct {
	Author author `json:"author"`
}
type author struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	
	Contributions []struct {
		Book hardcoverBook `json:"book"`
	} `json:"contributions"`
}

type bookSeries struct {
	ID       int     `json:"id"`
	Position float64 `json:"position"`
	Series   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"series"`
}

type hardcoverBook struct {
	ID                   int           `json:"id"`
	Title                string        `json:"title"`
	Contributors         []contributor `json:"contributions"`
	UsersCount           int           `json:"users_count"`
	FeaturedBookSeriesID int           `json:"featured_book_series_id"`
	BookSeries           []bookSeries  `json:"book_series"`
	Slug                 string        `json:"slug"`
	Subtitle             string        `json:"subtitle"`
}

type userBook struct {
	Book hardcoverBook `json:"book"`
}

type wantToReadResponse struct {
	Data struct {
		Me []struct {
			UserBooks []userBook `json:"user_books"`
		} `json:"me"`
	} `json:"data"`
}

type authorResponse struct {
	Data struct {
		FilterLists []struct {
			AuthorsCount       int `json:"authors_count"`
			FilterListEntities []struct {
				Author author `json:"author"`
			} `json:"filter_list_entities"`
		} `json:"filter_lists"`
	} `json:"data"`
}


