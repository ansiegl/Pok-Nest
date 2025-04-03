package collection_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCollectionPokemonDetailSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon/"+fixtures.PokemonInCollection1.PokemonID, nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		ctx := context.Background()
		exists, err := models.CollectionPokemons(
			models.CollectionPokemonWhere.UserID.EQ(fixtures.User1.ID),
			models.CollectionPokemonWhere.PokemonID.EQ(fixtures.PokemonInCollection1.PokemonID),
		).Exists(ctx, s.DB)
		require.NoError(t, err)
		assert.True(t, exists, "Das Pokémon sollte in der Sammlung existieren")
	})
}

func TestGetCollectionPokemonDetailNotFound(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		nonExistentPokemonID := "non-existent-id"

		res := test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon/"+nonExistentPokemonID, nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	})
}

func TestGetCollectionPokemonDetailUnauthorized(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "GET", "/api/v1/collection/pokemon/"+fixtures.PokemonInCollection1.PokemonID, nil, nil)
		require.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)
	})
}
