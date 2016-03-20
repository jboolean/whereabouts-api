package api

type latLngResponse struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type rawLocationResponse struct {
	Position latLngResponse `json:"position"`
}
