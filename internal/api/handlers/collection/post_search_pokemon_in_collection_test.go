package collection_test

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

func TestPostSearchPokemonInCollectio(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "POST", "/api/v1/collection/pokemon", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
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

func TestPostSearchPokemonInCollectionWithBody(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		// filter
		searchRequest := map[string]interface{}{
			"type": fixtures.PokemonInCollection1.Type1,
			"hp":   fixtures.PokemonInCollection1.HP,
			"pagination": map[string]interface{}{
				"limit":  2,
				"offset": 0,
			},
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/collection/pokemon", searchRequest, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// check pokemon in response
		var response types.GetPokemonResponse
		err := json.NewDecoder(res.Result().Body).Decode(&response)
		require.NoError(t, err)

		// check pokemon responses
		assert.Equal(t, fixtures.PokemonInCollection1.Type1, *response.Data[0].Type1, "Types doesn't match")
		assert.Equal(t, fixtures.PokemonInCollection1.HP, int(response.Data[0].Hp), "HP doesn#t match")

		// check pagination
		assert.Equal(t, 2, int(response.Pagination.Limit), "Limit should be 2")
		assert.Equal(t, 0, int(response.Pagination.Offset), "Offset should be 0")
	})
}
