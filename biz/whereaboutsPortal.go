package biz

import (
	"github.com/jboolean/whereabouts-api/dao"
	"github.com/jboolean/whereabouts-api/model"
	"log"
	"time"
)

// To get a smooth velocity we average the change from all the recorded locations in the window
//
// why not 0? Data 2 seconds apart does not tell us velocity very well.
const VelocityLookbackWindowStart = 2 * time.Minute
const VelocityLookbackWindowEnd = 17 * time.Minute

func ProcessAndStoreLocation(rawLocation model.LatLng, username string) {
	duration, distance := CalculateDistanceFromHome(rawLocation)
	processedLocation := &model.ProcessedLocation{
		TimeToHome:       duration,
		DistanceFromHome: distance,
		Username:         username,
	}

	log.Printf("Storing new location %#v", processedLocation)

	dao.ProcessedLocationsDAO.Store(processedLocation)

	go CreateAndStoreSummary(username)
}

func CreateAndStoreSummary(username string) {
	latestLocation := dao.ProcessedLocationsDAO.FindLatest(username)
	if latestLocation == nil {
		return
	}

	// calculate velocity
	recentLocations := dao.ProcessedLocationsDAO.FindInTimeRange(
		username,
		latestLocation.CreatedOn.Add(-VelocityLookbackWindowStart),
		latestLocation.CreatedOn.Add(-VelocityLookbackWindowEnd))

	sum := 0.0

	for _, loc := range recentLocations {
		dx := latestLocation.DistanceFromHome - loc.DistanceFromHome
		dt := latestLocation.CreatedOn.Unix() - loc.CreatedOn.Unix()
		sum += float64(dx) / float64(dt)
	}
	var averageVelocity = 0.0
	if len(recentLocations) > 0 {
		averageVelocity = sum / float64(len(recentLocations))
	}

	latestSummary := &model.WhereaboutsSummary{
		Username:   username,
		UpdatedOn:  latestLocation.CreatedOn,
		TimeToHome: latestLocation.TimeToHome,
		Velocity:   averageVelocity,
	}

	dao.WhereaboutsSummaryDAO.Store(latestSummary)
}

func FindAllSummaries() []model.WhereaboutsSummary {
	return dao.WhereaboutsSummaryDAO.FindAll()
}

func FindSummaryByUsername(username string) *model.WhereaboutsSummary {
	return dao.WhereaboutsSummaryDAO.FindByUsername(username)
}
