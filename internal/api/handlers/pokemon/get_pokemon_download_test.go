package pokemon_test

import (
	"context"
	"encoding/csv"
	"net/http"
	"strconv"
	"testing"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPokemonDownload(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "GET", "/api/v1/pokemon/download", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// check headers
		require.Equal(t, "text/csv", res.Result().Header.Get("Content-Type"))
		require.Equal(t, "attachment; filename=pokemon.csv", res.Result().Header.Get("Content-Disposition"))

		reader := csv.NewReader(res.Result().Body)
		csvData, err := reader.ReadAll()
		require.NoError(t, err)

		totalCount, err := models.Pokemons().Count(ctx, s.DB)
		require.NoError(t, err)

		// check length and first row (title) of csvData
		assert.Equal(t, int(totalCount)+1, len(csvData))
		assert.Equal(t, []string{"PokemonID", "Number", "Name", "Type 1", "Type 2", "HP", "Attack", "Defense", "Speed", "Special", "Gif URL", "PNG URL", "Description"}, csvData[0])

		// check first pokemon
		assert.Equal(t, []string{fixtures.PokemonNotInCollection.PokemonID,
			strconv.Itoa(fixtures.PokemonNotInCollection.PokemonNumber),
			fixtures.PokemonNotInCollection.Name,
			fixtures.PokemonNotInCollection.Type1,
			fixtures.PokemonNotInCollection.Type2.String,
			strconv.Itoa(fixtures.PokemonNotInCollection.HP),
			strconv.Itoa(fixtures.PokemonNotInCollection.Attack),
			strconv.Itoa(fixtures.PokemonNotInCollection.Defense),
			strconv.Itoa(fixtures.PokemonNotInCollection.Speed),
			strconv.Itoa(fixtures.PokemonNotInCollection.Special),
			fixtures.PokemonNotInCollection.GifURL,
			fixtures.PokemonNotInCollection.PNGURL,
			fixtures.PokemonNotInCollection.Description,
		}, csvData[1])

	})

}
