package collection

import (
	"fmt"
	"net/http"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/auth"
	"github.com/ansiegl/Pok-Nest.git/internal/api/httperrors"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/ansiegl/Pok-Nest.git/internal/types/collection"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetPokemonDetailsRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Collection.GET("/pokemon/:pokemonId", getPokemonDetailsHandler(s))
}

func getPokemonDetailsHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)
		user := auth.UserFromContext(ctx)

		params := collection.NewGetPokemonInCollectionDetailsParams()
		err := util.BindAndValidatePathAndQueryParams(c, &params)
		if err != nil {
			return err
		}

		queryMods := []qm.QueryMod{
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
			qm.InnerJoin("pokemon ON pokemon.pokemon_id = collection_pokemon.pokemon_id"),
			qm.Load("Pokemon"),
		}

		queryMods = append(queryMods, qm.Where("pokemon.pokemon_id = ?", params.PokemonID))

		collectionPokemon, err := models.CollectionPokemons(queryMods...).One(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Pokemon not found")
			return httperrors.NewHTTPError(http.StatusNotFound, "POKEMON_NOT_FOUND", fmt.Sprintf("No pokemon found that matches the ID %v", params.PokemonID))
		}

		pokemon := collectionPokemon.R.Pokemon
		if pokemon == nil {
			return httperrors.NewHTTPError(http.StatusNotFound, "POKEMONDETAILS_NOT_FOUND", fmt.Sprintf("Pokemon details for %v not available", params.PokemonID))
		}

		pokemonIDStr := strfmt.UUID4(pokemon.PokemonID)
		response := &types.Pokemon{
			PokemonID:   &pokemonIDStr,
			Name:        &pokemon.Name,
			Type1:       &pokemon.Type1,
			Type2:       pokemon.Type2.String,
			Hp:          int64(pokemon.HP),
			Attack:      int64(pokemon.Attack),
			Defense:     int64(pokemon.Defense),
			Speed:       int64(pokemon.Speed),
			Special:     int64(pokemon.Special),
			Png:         pokemon.PNGURL,
			Description: pokemon.Description,
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
