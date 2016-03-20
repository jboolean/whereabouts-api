package api

import (
	"github.com/gorilla/mux"
	"github.com/jboolean/whereabouts-api/biz"
	"github.com/jboolean/whereabouts-api/model"
	"github.com/jboolean/whereabouts-api/router"
	"net/http"
)

var userRoutes = router.Routes{
	router.Route{
		Name:    "CreateUser",
		Method:  "POST",
		Pattern: "/users",
		Handler: requiresAdmin(jsonResource(createUser)),
	},
	router.Route{
		Name:    "RetrieveUser",
		Method:  "GET",
		Pattern: "/users/{username}",
		Handler: requiresRead(jsonResource(getUser)),
	},
	router.Route{
		Name:    "DeleteUser",
		Method:  "DELETE",
		Pattern: "/users/{username}",
		Handler: isUserOrAdmin(jsonResource(deleteUser)),
	},
	router.Route{
		Name:    "ChangePassword",
		Method:  "PUT",
		Pattern: "/users/{username}/password",
		Handler: isUserOrAdmin(jsonResource(changePassword)),
	},
}

func createUser(r *http.Request) (interface{}, *resourceError) {
	userRequest := new(userResponse)

	var err error
	err = decodeJsonBody(userRequest, r)
	if err != nil {
		return nil, &resourceError{err, "Error reading request body.", http.StatusBadRequest}
	}

	user := &model.User{
		Username:    userRequest.Username,
		DisplayName: userRequest.DisplayName,
	}
	// todo password

	err = biz.CreateUser(user, userRequest.Roles, userRequest.Password)
	if err != nil {
		if _, ok := err.(biz.DuplicateUsernameError); ok {
			return nil, &resourceError{err, "Username taken", http.StatusBadRequest}
		}
		return nil, &resourceError{err, "Error creating user", http.StatusInternalServerError}
	}

	var roles []model.Role

	user, err = biz.FindUserByUsername(user.Username)
	roles, err = biz.FindRolesByUsername(user.Username)

	if err != nil {
		return nil, &resourceError{err, "Error retrieving user", http.StatusInternalServerError}
	}

	// TODO add roles to request and response
	return makeUserResponse(user, roles), nil
}

func getUser(r *http.Request) (interface{}, *resourceError) {
	vars := mux.Vars(r)
	username := vars["username"]
	var user *model.User
	var roles []model.Role
	var err error

	user, err = biz.FindUserByUsername(username)
	roles, err = biz.FindRolesByUsername(username)

	if err != nil {
		return nil, &resourceError{err, "Error retrieving user", http.StatusInternalServerError}
	}
	if user == nil {
		return nil, &resourceError{nil, "User not found", http.StatusNotFound}
	}

	return makeUserResponse(user, roles), nil
}

func deleteUser(r *http.Request) (interface{}, *resourceError) {
	vars := mux.Vars(r)
	username := vars["username"]
	var user *model.User

	user, _ = biz.FindUserByUsername(username)
	if user == nil {
		return nil, &resourceError{nil, "User not found", http.StatusNotFound}
	}

	biz.DeleteUser(username)
	return nil, nil
}

func changePassword(r *http.Request) (interface{}, *resourceError) {
	vars := mux.Vars(r)
	username := vars["username"]
	var user *model.User

	user, _ = biz.FindUserByUsername(username)
	if user == nil {
		return nil, &resourceError{nil, "User not found", http.StatusNotFound}
	}

	password := readBodyString(r)
	biz.UpdateUserPassword(username, password)
	return nil, nil
}
