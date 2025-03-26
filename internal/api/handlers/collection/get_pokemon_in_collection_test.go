package collection_test

import (
	"net/http"
	"testing"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGetPokemonInCollection(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		res = test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon?type=fire", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		res = test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon?generation=1", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		res = test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon?legendary=true", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		res = test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon?type=fire&generation=1&legendary=true", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		res = test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon?limit=5&offset=0", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}
