package biz

import (
	"github.com/jboolean/whereabouts-api/model"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var mapClient *maps.Client

var travelMode maps.Mode
var Home model.LatLng

func init() {
	var err error
	var mapsApiKey = os.Getenv("WA_MAPS_API_KEY")
	travelMode = maps.Mode(os.Getenv("WA_TRAVEL_MODE"))

	homeString := os.Getenv("WA_HOME")
	splitHomeString := strings.Split(homeString, ",")
	if len(splitHomeString) < 2 {
		log.Fatalf("Environment variable WA_HOME must be in the format 000.00,000.00. Found %s", homeString)
	}

	homeLat, err := strconv.ParseFloat(splitHomeString[0], 64)
	homeLng, err := strconv.ParseFloat(splitHomeString[1], 64)
	if err != nil {
		log.Fatalf("Environment variable WA_HOME must be in the format 000.00,000.00. Found %s", homeString)
	}
	Home = model.LatLng{homeLat, homeLng}

	mapClient, err = maps.NewClient(maps.WithAPIKey(mapsApiKey))
	if err != nil {
		log.Fatalf("Could not start Google Maps client. %v", err)
	}
}

// Returns distance between latLng and home in time and meters
func CalculateDistanceFromHome(latLng model.LatLng) (time.Duration, int) {
	request := &maps.DistanceMatrixRequest{
		Origins:      []string{Home.String()},
		Destinations: []string{latLng.String()},
		Mode:         travelMode,
	}

	response, err := mapClient.DistanceMatrix(context.Background(), request)
	if err != nil {
		log.Panic(err)
	}
	cell := response.Rows[0].Elements[0]
	return cell.Duration, cell.Distance.Meters
}
