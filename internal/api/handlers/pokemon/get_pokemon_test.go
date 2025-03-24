package pokemon_test

import (
	"net/http"
	"testing"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGetPokemon(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		// test without filter
		res := test.PerformRequest(t, s, "GET", "/api/v1/pokemon", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// test fitler "type"
		res = test.PerformRequest(t, s, "GET", "/api/v1/pokemon?type=fire", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// test fitler "generation"
		res = test.PerformRequest(t, s, "GET", "/api/v1/pokemon?generation=1", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// test fitler "legendary"
		res = test.PerformRequest(t, s, "GET", "/api/v1/pokemon?legendary=true", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// test multiple filters: "type", "generation", "legendary"
		res = test.PerformRequest(t, s, "GET", "/api/v1/pokemon?type=fire&generation=1&legendary=true", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// test "limit" and "offset"
		res = test.PerformRequest(t, s, "GET", "/api/v1/pokemon?limit=5&offset=0", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}
