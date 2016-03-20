package api

import "github.com/jboolean/whereabouts-api/model"

type whereaboutsSummaryResponse struct {
	Username   string  `json:"username"`
	UpdatedOn  int64   `json:"updatedOn"`
	TimeToHome int64   `json:"timeToHome"`
	Velocity   float64 `json:"velocity"`
}

func makeWhereaboutsSummaryResponse(input *model.WhereaboutsSummary) *whereaboutsSummaryResponse {
	return &whereaboutsSummaryResponse{
		input.Username, input.UpdatedOn.Unix(), int64(input.TimeToHome.Seconds()), input.Velocity,
	}
}
