package collection_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/stretchr/testify/require"
)

func TestDeletePokemon(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		// check if pokemon exists in collection before deleting
		existsBefore, err := models.CollectionPokemons(
			models.CollectionPokemonWhere.PokemonID.EQ(fixtures.PokemonInCollection1.PokemonID),
			models.CollectionPokemonWhere.UserID.EQ(fixtures.User1.ID),
		).Exists(context.Background(), s.DB)
		require.NoError(t, err)
		require.True(t, existsBefore, "Pokemon sollte vor dem Löschen in der Sammlung existieren")

		// delete pokemon
		res := test.PerformRequest(t, s, "DELETE", "/api/v1/collection/pokemon/"+fixtures.PokemonInCollection1.PokemonID, nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// After deleting, check if pokemon exists in collection
		existsAfter, err := models.CollectionPokemons(
			models.CollectionPokemonWhere.PokemonID.EQ(fixtures.PokemonInCollection1.PokemonID),
			models.CollectionPokemonWhere.UserID.EQ(fixtures.User1.ID),
		).Exists(context.Background(), s.DB)
		require.NoError(t, err)
		require.False(t, existsAfter, "Pokemon sollte nach dem Löschen nicht mehr in der Sammlung existieren")
	})
}
