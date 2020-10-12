package server

import (
	"net/http"
	"testing"
)

func TestCommand_hasAdminToken(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost", nil)

	t.Run("has", func(in *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost", nil)
		request.Header.Set(adminTokenHeader, "this-is-a-token")

		if !hasAdminToken(request, "this-is-a-token") {
			in.Errorf("expected to be true, got false instead.")
		}
	})

	t.Run("has-not", func(in *testing.T) {
		request.Header.Del(adminTokenHeader)
		if hasAdminToken(request, "this-is-a-token") {
			in.Errorf("expected to be `false`, got `true` instead.")
		}
	})
}
