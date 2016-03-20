package api

import (
	"github.com/gorilla/mux"
	"github.com/jboolean/whereabouts-api/biz"
	"github.com/jboolean/whereabouts-api/model"
	"github.com/jboolean/whereabouts-api/router"
	"net/http"
)

var whereaboutsRoutes = router.Routes{
	router.Route{
		Name:    "SubmitLocation",
		Method:  "POST",
		Pattern: "/whereabouts/raw/{username}",
		Handler: isUserInPath(requiresMaintainLocation(jsonResource(submitLocation))),
	},
	router.Route{
		Name:    "GetCurrentWhereaboutsSummariesForAllUsers",
		Method:  "GET",
		Pattern: "/whereabouts/summaries/current",
		Handler: requiresRead(jsonResource(getCurrentWhereaboutsSummaryForAllUsers)),
	},
	router.Route{
		Name:    "GetCurrentWhereaboutsSummary",
		Method:  "GET",
		Pattern: "/whereabouts/summaries/{username}/current",
		Handler: requiresRead(jsonResource(getCurrentWhereaboutsSummary)),
	},
}

func submitLocation(r *http.Request) (interface{}, *resourceError) {
	vars := mux.Vars(r)
	username := vars["username"]

	submittedLocation := new(rawLocationResponse)

	var err error
	err = decodeJsonBody(submittedLocation, r)
	if err != nil {
		return nil, &resourceError{err, "Error reading request body.", http.StatusBadRequest}
	}

	userLocation := &model.LatLng{submittedLocation.Position.Lat, submittedLocation.Position.Lng}

	go biz.ProcessAndStoreLocation(*userLocation, username)

	return nil, nil
}

func getCurrentWhereaboutsSummaryForAllUsers(r *http.Request) (interface{}, *resourceError) {
	summaries := biz.FindAllSummaries()
	responseResults := make([]whereaboutsSummaryResponse, len(summaries))
	for i, v := range summaries {
		responseResults[i] = *makeWhereaboutsSummaryResponse(&v)
	}
	return responseResults, nil
}

func getCurrentWhereaboutsSummary(r *http.Request) (interface{}, *resourceError) {
	vars := mux.Vars(r)
	username := vars["username"]

	result := biz.FindSummaryByUsername(username)
	if result == nil {
		return nil, &resourceError{nil, "None found", http.StatusNotFound}
	}
	return makeWhereaboutsSummaryResponse(result), nil
}
