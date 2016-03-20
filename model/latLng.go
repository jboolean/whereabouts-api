package model

import "fmt"

type LatLng struct {
	Lat float64
	Lng float64
}

func (latLng LatLng) String() string {
	return fmt.Sprintf("%f,%f", latLng.Lat, latLng.Lng)
}
