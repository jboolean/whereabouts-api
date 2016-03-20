package api

import "github.com/jboolean/whereabouts-api/model"

type userResponse struct {
	Username    string       `json:"username"`
	DisplayName string       `json:"displayName"`
	Roles       []model.Role `json:"roles"`
	Password    string       `json:"password,omitempty"`
}

func makeUserResponse(user *model.User, roles []model.Role) *userResponse {
	return &userResponse{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Roles:       roles,
	}
}
