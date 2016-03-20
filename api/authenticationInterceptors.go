package api

import (
	"github.com/gorilla/mux"
	"github.com/jboolean/whereabouts-api/biz"
	"github.com/jboolean/whereabouts-api/model"
	"net/http"
)

// Wrap your handlers in these for protection :-)

type Interceptor func(http.Handler) http.Handler

var notAuthorizedError = resourceError{nil, "Not authorized.", http.StatusUnauthorized}
var notLoggedInError = resourceError{nil, "Not logged in", http.StatusUnauthorized}

func makeAuthenticationInterceptor(role model.Role) Interceptor {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if hasRoleHelper(r, role) {
				inner.ServeHTTP(w, r)
				return
			}
			notAuthorizedError.WriteToResponseAsJson(w)
		})
	}
}

// Some common auth interceptors
var requiresAdmin = makeAuthenticationInterceptor(model.RoleAdmin)
var requiresRead = makeAuthenticationInterceptor(model.RoleRead)
var requiresMaintainLocation = makeAuthenticationInterceptor(model.RoleMaintainLocation)

func isUserInPath(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isUserInPathHelper(r) {
			notAuthorizedError.WriteToResponseAsJson(w)
			return
		}
		inner.ServeHTTP(w, r)
	})
}

// Must be the user in the path or have the admin role
func isUserOrAdmin(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !(isUserInPathHelper(r) || hasRoleHelper(r, model.RoleAdmin)) {
			notAuthorizedError.WriteToResponseAsJson(w)
			return
		}
		inner.ServeHTTP(w, r)
	})
}

func hasRoleHelper(r *http.Request, role model.Role) bool {
	sessionKeyHeader := r.Header["X-Session-Key"]
	if len(sessionKeyHeader) < 1 {
		return false
	}
	sessionKey := sessionKeyHeader[0]

	session := biz.FindSessionByKey(sessionKey)
	if session == nil {
		return false
	}

	roles, _ := biz.FindRolesByUsername(session.Username)
	return roles.Contains(role)
}

func isUserInPathHelper(r *http.Request) bool {
	vars := mux.Vars(r)
	username := vars["username"]
	sessionKey := r.Header["X-Session-Key"][0]

	session := biz.FindSessionByKey(sessionKey)
	if session == nil {
		return false
	}

	return session.Username == username
}
