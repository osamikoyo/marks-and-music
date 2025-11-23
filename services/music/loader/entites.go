package loader

import "github.com/osamikoyo/music-and-marks/services/music/entity"

type Release struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status,omitempty"` // Official, Bootleg
	Quality   string `json:"quality,omitempty"`
	Date      string `json:"date,omitempty"` // YYYY-MM-DD
	Country   string `json:"country,omitempty"`
	Barcode   string `json:"barcode,omitempty"`
	LabelInfo []struct {
		Label struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"label"`
		CatalogNumber string `json:"catalog-number"`
	} `json:"label-info"`
	TrackCount int `json:"track-count,omitempty"`
	Media      []struct {
		Format     string `json:"format,omitempty"`
		TrackCount int    `json:"track-count"`
	} `json:"media"`
	ReleaseGroup struct {
		ID             string   `json:"id"`
		Title          string   `json:"title"`
		PrimaryType    string   `json:"primary-type"` // Album, Single, EP
		SecondaryTypes []string `json:"secondary-types,omitempty"`
	} `json:"release-group"`
	ArtistCredit []struct {
		Name       string `json:"name"`
		JoinPhrase string `json:"joinphrase,omitempty"`
		Artist     struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"artist"`
	} `json:"artist-credit"`
}

func (r *Release) ToEntity() *entity.Release {
	var date *string
	if r.Date != "" {
		date = &r.Date
	}

	var format string
	var trackCount int
	if len(r.Media) > 0 {
		format = r.Media[0].Format
		trackCount = r.Media[0].TrackCount
	}

	return &entity.Release{
		MBID:           r.ID,
		Title:          r.Title,
		ReleaseGroupID: r.ReleaseGroup.ID,
		Status:         r.Status,
		Country:        r.Country,
		Date:           date,
		Format:         format,
		TrackCount:     trackCount,
	}
}

type ReleaseSearchResult struct {
	Created  string    `json:"created"`
	Count    int       `json:"count"`
	Offset   int       `json:"offset"`
	Releases []Release `json:"releases"`
}

type Artist struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SortName  string `json:"sort-name"`
	Country   string `json:"country,omitempty"`
	Type      string `json:"type,omitempty"`
	TypeID    string `json:"type-id,omitempty"`
	Gender    string `json:"gender,omitempty"`
	GenderID  string `json:"gender-id,omitempty"`
	BeginDate string `json:"begin-date,omitempty"`
	EndDate   string `json:"end-date,omitempty"`
}

func (a *Artist) ToEntity() *entity.Artist {
	return &entity.Artist{
		ID:       a.ID,
		Name:     a.Name,
		Country:  a.Country,
		Type:     a.Type,
		SortName: a.SortName,
	}
}

type ArtistSearchResult struct {
	Created string   `json:"created"`
	Count   int      `json:"count"`
	Offset  int      `json:"offset"`
	Artists []Artist `json:"artists"`
}
