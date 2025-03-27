package collection_test

import (
	"net/http"
	"testing"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/stretchr/testify/require"
)

func TestDeletePokemon(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "DELETE", "/api/v1/collection/pokemon/"+fixtures.PokemonInCollection1.PokemonID, nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}
