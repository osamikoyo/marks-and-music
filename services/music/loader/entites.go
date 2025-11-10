package loader

type Release struct {
    ID             string   `json:"id"`
    Title          string   `json:"title"`
    Status         string   `json:"status,omitempty"`         // Official, Bootleg
    Quality        string   `json:"quality,omitempty"`
    Date           string   `json:"date,omitempty"`           // YYYY-MM-DD
    Country        string   `json:"country,omitempty"`
    Barcode        string   `json:"barcode,omitempty"`
    LabelInfo      []struct {
        Label struct {
            ID   string `json:"id"`
            Name string `json:"name"`
        } `json:"label"`
        CatalogNumber string `json:"catalog-number"`
    } `json:"label-info"`
    TrackCount     int `json:"track-count,omitempty"`
    Media          []struct {
        Format     string `json:"format,omitempty"`
        TrackCount int    `json:"track-count"`
    } `json:"media"`
    ReleaseGroup struct {
        ID         string `json:"id"`
        Title      string `json:"title"`
        PrimaryType string `json:"primary-type"` // Album, Single, EP
        SecondaryTypes []string `json:"secondary-types,omitempty"`
    } `json:"release-group"`
    ArtistCredit []struct {
        Name     string `json:"name"`
        JoinPhrase string `json:"joinphrase,omitempty"`
        Artist   struct {
            ID   string `json:"id"`
            Name string `json:"name"`
        } `json:"artist"`
    } `json:"artist-credit"`
}

type ReleaseSearchResult struct {
    Created   string    `json:"created"`
    Count     int       `json:"count"`
    Offset    int       `json:"offset"`
    Releases  []Release `json:"releases"`
}

type Artist struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    SortName    string `json:"sort-name"`
    Country     string `json:"country,omitempty"`
    Type        string `json:"type,omitempty"`
    TypeID      string `json:"type-id,omitempty"`
    Gender      string `json:"gender,omitempty"`
    GenderID    string `json:"gender-id,omitempty"`
    BeginDate   string `json:"begin-date,omitempty"`
    EndDate     string `json:"end-date,omitempty"`
}

type SearchResult struct {
    Created string    `json:"created"`
    Count   int       `json:"count"`
    Offset  int       `json:"offset"`
    Artists []Artist  `json:"artists"`
}