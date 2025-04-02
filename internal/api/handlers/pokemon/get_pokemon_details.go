package pokemon

import (
	"net/http"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/httperrors"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/ansiegl/Pok-Nest.git/internal/types/collection"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

func GetPokemonDetailsRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Pokemon.GET("/:pokemonId", getPokemonDetailsHandler(s))
}

func getPokemonDetailsHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		params := collection.NewPostAddPokemonToCollectionParams()
		err := util.BindAndValidatePathAndQueryParams(c, &params)
		if err != nil {
			return err
		}

		pokemonExists, err := models.Pokemons(models.PokemonWhere.PokemonID.EQ(params.PokemonID)).Exists(ctx, s.DB)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to check if pokemon exsists")
		}
		if !pokemonExists {
			log.Debug().Str("pokemonID", params.PokemonID).Msg("Pokemon does not exist")
			return httperrors.NewHTTPError(http.StatusNotFound, "POKEMON_NOT_FOUND", "Pokemon not found")
		}

		pokemon, err := models.Pokemons(models.PokemonWhere.PokemonID.EQ(params.PokemonID)).One(ctx, s.DB)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to get pokemon")
			return httperrors.NewHTTPError(http.StatusInternalServerError, "POKEMON_NOT_FOUND", "Pokemon not found")
		}

		pokemonID := strfmt.UUID4(pokemon.PokemonID)
		response := &types.Pokemon{
			PokemonID:   &pokemonID,
			Name:        &pokemon.Name,
			Type1:       &pokemon.Type1,
			Type2:       pokemon.Type2.String,
			Hp:          int64(pokemon.HP),
			Attack:      int64(pokemon.Attack),
			Defense:     int64(pokemon.Defense),
			Speed:       int64(pokemon.Speed),
			Special:     int64(pokemon.Special),
			Gif:         pokemon.GifURL,
			Png:         pokemon.PNGURL,
			Description: pokemon.Description,
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)

	}
}
