package collection_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/stretchr/testify/require"
)

func TestGetPokemon(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		// test GET all pokemon with default pagination
		res := test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// parse response body
		var response types.GetPokemonResponse
		err := json.NewDecoder(res.Result().Body).Decode(&response)
		require.NoError(t, err)

		// check that pagination is correct
		require.NotNil(t, response.Pagination)
		require.NotNil(t, response.Pagination.Total)
		require.NotNil(t, response.Pagination.Limit)
		require.NotNil(t, response.Pagination.Offset)
		require.Equal(t, int64(10), response.Pagination.Limit) // Default limit is 10
		require.Equal(t, int64(0), response.Pagination.Offset) // Default offset is 0

		// check that data was returned
		require.NotNil(t, response.Data)
		require.NotEmpty(t, response.Data)

		// test limit parameter
		res = test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon?limit=3", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		var limitResponse types.GetPokemonResponse
		err = json.NewDecoder(res.Result().Body).Decode(&limitResponse)
		require.NoError(t, err)

		// check that only 3 pokemon were returned
		require.Len(t, limitResponse.Data, 3)
		require.Equal(t, int64(3), limitResponse.Pagination.Limit)
	})
}
