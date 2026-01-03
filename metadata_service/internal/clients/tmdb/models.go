package tmdb

type TMDBMovieResponse struct {
	Results    []TMDBMovie `json:"results"`
	TotalPages int         `json:"total_pages"`
}

type TMDBMovie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	PosterPath  *string `json:"poster_path"`
	GenreIDs    []int   `json:"genre_ids"`
	Overview    string  `json:"overview"`
}

type TMDBTVResponse struct {
	Results    []TMDBTV `json:"results"`
	TotalPages int      `json:"total_pages"`
}

type TMDBTV struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	FirstAirDate string  `json:"first_air_date"`
	PosterPath   *string `json:"poster_path"`
	GenreIDs     []int   `json:"genre_ids"`
	Overview     string  `json:"overview"`
}

type TMDBMovieDetails struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	PosterPath  *string `json:"poster_path"`
	Genres      []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Overview string `json:"overview"`
}

type TMDBTVDetails struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	FirstAirDate string  `json:"first_air_date"`
	PosterPath   *string `json:"poster_path"`
	Genres       []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Overview string `json:"overview"`
}

