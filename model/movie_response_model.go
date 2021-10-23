package model

type Movies struct {
	Search []Movie `json:"Search"`
	Response string `json:"Response"`
	Error string `json:"Error"`
}

type MovieDetailResponse struct {
	MovieDetail MovieDetail `json:"movie_detail"`
	Error string
}

type MovieDetail struct {
	Title string
	Year string
	Rated string
	Released string
	Runtime string
	Genre string
	Director string
	Writer string
	Actors string
	Plot string
	Language string
	Country string
	Awards string
	Poster string
	Ratings []Rating
	Metasource string
	imdbRating string
	imdbVotes string
	imdbID string
	Type string
	DVD string
	BoxOffice string
	Production string
	Website string
	Response string
}

type Rating struct {
	Source string
	Value string
}

type Movie struct {
	ImdbID string `json:"imdbID"`
	Title  string `json:"Title"`
	Year   string    `json:"Year"`
	Type   string `json:"Type"`
	Poster string `json:"Poster"`
}
