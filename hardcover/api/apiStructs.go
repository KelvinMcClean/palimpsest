package api

type contributor struct {
	Author author `json:"author"`
}
type author struct {
	Name string `json:"name"`
}

type hardcoverBook struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	Contributors []contributor `json:"contributions"`
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

type Book struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
}
