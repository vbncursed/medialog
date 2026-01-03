package open_library

type OpenLibrarySearchResponse struct {
	Docs     []OpenLibraryBook `json:"docs"`
	NumFound int               `json:"numFound"`
}

type OpenLibraryBook struct {
	Key              string   `json:"key"`
	Title            string   `json:"title"`
	FirstPublishYear *int     `json:"first_publish_year"`
	CoverI           *int     `json:"cover_i"`
	ISBN             []string `json:"isbn"`
	AuthorName       []string `json:"author_name"`
	Subject          []string `json:"subject"`
}

type OpenLibraryBookDetails struct {
	Key              string   `json:"key"`
	Title            string   `json:"title"`
	FirstPublishYear *int     `json:"first_publish_year"`
	Covers           []int    `json:"covers"`
	ISBN10           []string `json:"isbn_10"`
	ISBN13           []string `json:"isbn_13"`
	Authors          []struct {
		Key string `json:"key"`
	} `json:"authors"`
	Subjects []string `json:"subjects"`
}
