package collection_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPutEditPokemon(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		ctx := context.Background()

		// get pokemon before update
		originalPokemon, err := models.CollectionPokemons(
			models.CollectionPokemonWhere.PokemonID.EQ(fixtures.PokemonInCollection1.PokemonID),
			models.CollectionPokemonWhere.UserID.EQ(fixtures.User1.ID),
		).One(ctx, s.DB)
		require.NoError(t, err)

		// set new values
		newNickname := "TestNickname"
		newCaughtTime := time.Date(0, 1, 1, 15, 9, 0, 0, time.UTC)

		updateBody := map[string]interface{}{
			"nickname": newNickname,
			"caught":   newCaughtTime,
		}

		res := test.PerformRequest(t, s, "PUT",
			"/api/v1/collection/pokemon/"+fixtures.PokemonInCollection1.PokemonID,
			updateBody,
			test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// check pokemon after update
		updatedPokemon, err := models.CollectionPokemons(
			models.CollectionPokemonWhere.PokemonID.EQ(fixtures.PokemonInCollection1.PokemonID),
			models.CollectionPokemonWhere.UserID.EQ(fixtures.User1.ID),
		).One(ctx, s.DB)
		require.NoError(t, err)

		assert.NotEqual(t, originalPokemon.Nickname.String, updatedPokemon.Nickname.String,
			"Nickname didn't change")
		assert.Equal(t, newNickname, updatedPokemon.Nickname.String,
			"Nickname updated")
	})
}
