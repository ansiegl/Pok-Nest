package collection

import (
	"net/http"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/auth"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
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

		pokemonID := c.Param("pokemonId")
		if pokemonID == "" {
			log.Debug().Msg("Pokemon ID is missing from path")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Pokemon ID is missing"})
		}

		exists, err := models.CollectionPokemons(
			models.CollectionPokemonWhere.PokemonID.EQ(pokemonID),
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
		).Exists(ctx, s.DB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check if pokemon exists in collection")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}

		if !exists {
			log.Debug().Str("pokemonID", pokemonID).Msg("Pokemon does not exist in collection")
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Pokemon not found in collection"})
		}

		_, err = models.CollectionPokemons(
			models.CollectionPokemonWhere.PokemonID.EQ(pokemonID),
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
		).DeleteAll(ctx, s.DB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete pokemon from collection")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete pokemon"})
		}

		log.Info().Str("pokemonID", pokemonID).Str("userID", user.ID).Msg("Pokemon successfully deleted from collection")
		return c.JSON(http.StatusOK, map[string]string{"message": "Pokemon successfully deleted from collection"})
	}
}
