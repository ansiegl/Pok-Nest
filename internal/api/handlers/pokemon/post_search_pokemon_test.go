package pokemon_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostSearchPokemon(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "POST", "/api/v1/pokemon", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
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
		res = test.PerformRequest(t, s, "GET", "/api/v1/pokemon?limit=3", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		var limitResponse types.GetPokemonResponse
		err = json.NewDecoder(res.Result().Body).Decode(&limitResponse)
		require.NoError(t, err)

		// check that only 3 pokemon were returned
		require.Len(t, limitResponse.Data, 3)
		require.Equal(t, int64(3), limitResponse.Pagination.Limit)
	})
}

func TestPostSearchPokemonWithBody(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		// Test-Filter
		searchRequest := map[string]interface{}{
			"name": "Bulbasaur",
			"type": "Grass",
			"hp":   45,
			"pagination": map[string]interface{}{
				"limit":  10,
				"offset": 0,
			},
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/pokemon", searchRequest, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// check pokemon in response
		var response types.GetPokemonResponse
		err := json.NewDecoder(res.Result().Body).Decode(&response)
		require.NoError(t, err)

		// Überprüfen, dass das Pokémon in der Antwort enthalten ist
		assert.Len(t, response.Data, 1, "Expected one pokemon in the response")
		assert.Equal(t, "Bulbasaur", *response.Data[0].Name, "Name of pokemon should match")
		assert.Equal(t, "Grass", *response.Data[0].Type1, "Type of pokemon should match")
		assert.Equal(t, 45, int(response.Data[0].Hp), "HP of pokemon should match")

		// check pagination
		assert.Equal(t, 1, int(response.Pagination.Total), "Total should be 1")
		assert.Equal(t, 10, int(response.Pagination.Limit), "Limit should be 10")
		assert.Equal(t, 0, int(response.Pagination.Offset), "Offset should be 0")
	})
}
