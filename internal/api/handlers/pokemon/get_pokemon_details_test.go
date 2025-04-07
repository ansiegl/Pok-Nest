package pokemon_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/stretchr/testify/require"
)

func TestGetPokemonDetails(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		pokemonID := string(fixtures.PokemonInCollection1.PokemonID)

		// test valid request
		res := test.PerformRequest(t, s, "GET", "/api/v1/pokemon/"+pokemonID, nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.Pokemon
		err := json.NewDecoder(res.Result().Body).Decode(&response)
		require.NoError(t, err)

		// check returned Pokemon data
		require.Equal(t, pokemonID, string(*response.PokemonID))
		require.Equal(t, fixtures.PokemonInCollection1.Name, *response.Name)
		require.Equal(t, fixtures.PokemonInCollection1.Type1, *response.Type1)
		require.Equal(t, fixtures.PokemonInCollection1.PNGURL, response.ImageURL)

		// test invalid Pokemon ID
		invalidID := "00000000-0000-0000-0000-000000000000"
		res = test.PerformRequest(t, s, "GET", "/api/v1/pokemon/"+invalidID, nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	})
}
