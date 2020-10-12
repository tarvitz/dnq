package server

import "net/http"

const (
	adminTokenHeader = "X-Admin-Token"
)

// very basic and simple `auth`, however, it's discouraged to
// develop like this way and concentrate on standards way of authentication
// and authorization.
func hasAdminToken(request *http.Request, token string) bool {
	adminToken := request.Header.Get(adminTokenHeader)
	return adminToken == token
}
