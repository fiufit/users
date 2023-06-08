package utils

import (
	"strings"

	"github.com/sams96/rgeo"
)

type ReverseLocator struct {
	rGeo *rgeo.Rgeo
}

func NewReverseLocator() (*ReverseLocator, error) {
	r, err := rgeo.New(rgeo.Countries110)
	if err != nil {
		return &ReverseLocator{}, err
	}

	return &ReverseLocator{rGeo: r}, nil
}

func (rl ReverseLocator) GetLocationFromCoordinates(lat float64, long float64) (string, error) {
	loc, err := rl.rGeo.ReverseGeocode([]float64{long, lat})
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(loc.String(), "<Location> ", ""), nil
}
