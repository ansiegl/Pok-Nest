package collection_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/test"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
)

func TestPutEditPokemon(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		body := test.GenericPayload{
			"Caught":   null.TimeFrom(time.Now()),
			"Nickname": null.StringFrom("Das ist der beste Nickname der Welt"),
		}

		res := test.PerformRequest(t, s, "PUT", "/api/v1/collection/pokemon/"+fixtures.PokemonInCollection1.PokemonID, body, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}
