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

func TestPostAddPokemonToCollection(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "POST", "/api/v1/pokemon/"+fixtures.PokemonNotInCollection.PokemonID, nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusCreated, res.Result().StatusCode)

		// check if pokemon is in the collection
		ctx := context.Background()
		exists, err := models.CollectionPokemons(
			models.CollectionPokemonWhere.UserID.EQ(fixtures.User1.ID),
			models.CollectionPokemonWhere.PokemonID.EQ(fixtures.PokemonNotInCollection.PokemonID),
		).Exists(ctx, s.DB)
		require.NoError(t, err)
		assert.True(t, exists, "Pokemon should be added to the collection in the database")
	})
}

func TestPostAddPokemonToCollectionAlreadyExists(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "POST", "/api/v1/pokemon/"+fixtures.PokemonInCollection1.PokemonID, nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	})
}

func TestPostAddPokemonToCollectionNotFound(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		// fake id
		nonExistentPokemonID := "non-existent-id"

		res := test.PerformRequest(t, s, "POST", "/api/v1/pokemon/"+nonExistentPokemonID, nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	})
}

func TestPostAddPokemonToCollectionInvalidInput(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		reqBody := map[string]interface{}{
			"caught": "invalid-date-format",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/pokemon/"+fixtures.PokemonInCollection1.PokemonID, reqBody, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	})
}
