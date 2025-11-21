// Package entity stores entitites
package entity

import "github.com/osamikoyo/music-and-marks/services/music/api/proto/gen/pb"

type SearchResult struct {
	ID          string  `json:"id"`
	MBID        string  `json:"mbid,omitempty"`
	Title       string  `json:"title"`
	ArtistName  string  `json:"artist_name,omitempty"`
	Type        string  `json:"type"`
	ReleaseDate *string `json:"release_date,omitempty"`
	Relevance   float32 `json:"relevance"`
}

func (sr *SearchResult) ToPB() *pb.SearchResult {
	return &pb.SearchResult{
		Id:          sr.ID,
		Mbid:        sr.MBID,
		Title:       sr.Title,
		ArtistName:  sr.ArtistName,
		Type:        sr.Type,
		ReleaseDate: sr.ReleaseDate,
		Relevance:   sr.Relevance,
	}
}
