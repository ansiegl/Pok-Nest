package collection

import (
	"fmt"
	"net/http"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/auth"
	"github.com/ansiegl/Pok-Nest.git/internal/api/httperrors"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/types/collection"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/labstack/echo/v4"
)

func DeletePokemonRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Collection.DELETE("/pokemon/:pokemonId", deletePokemonHandler(s))
}

func deletePokemonHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		user := auth.UserFromContext(ctx)

		// get parameter
		params := collection.NewDeletePokemonFromCollectionParams()
		if err := util.BindAndValidatePathAndQueryParams(c, &params); err != nil {
			return err
		}

		exists, err := models.CollectionPokemons(
			models.CollectionPokemonWhere.PokemonID.EQ(params.PokemonID),
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
		).Exists(ctx, s.DB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check if pokemon exists in collection")
			return err
		}

		if !exists {
			log.Debug().Str("pokemonID", params.PokemonID).Msg("Pokemon does not exist in collection")
			return httperrors.NewHTTPError(http.StatusNotFound, "POKEMON_NOT_FOUND", fmt.Sprintf("No pokemon found that matches the ID %v", params.PokemonID))
		}

		_, err = models.CollectionPokemons(
			models.CollectionPokemonWhere.PokemonID.EQ(params.PokemonID),
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
		).DeleteAll(ctx, s.DB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete pokemon from collection")
			return err
		}

		log.Info().Str("pokemonID", params.PokemonID).Str("userID", user.ID).Msg("Pokemon successfully deleted from collection")
		return c.JSON(http.StatusOK, map[string]string{"message": "Pokemon successfully deleted from collection"})
	}
}
