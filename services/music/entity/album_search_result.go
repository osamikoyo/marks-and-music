package entity

type AlbumSearchResult struct {
    ID              string  `json:"id"`
    MBID            string  `json:"mbid"`
    Title           string  `json:"title"`
    ArtistName      string  `json:"artist_name"`
    FirstReleaseDate *string `json:"first_release_date"`
    AvgRating       float64 `json:"avg_rating"`
    ReviewCount     int64   `json:"review_count"`
    Relevance       float32 `json:"relevance"`
}