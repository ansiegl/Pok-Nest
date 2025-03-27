package collection

import (
	"net/http"
	"time"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/auth"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func PutEditPokemonRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Collection.PUT("/pokemon/:pokemonId", putEditPokemonHandler(s))
}

func putEditPokemonHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		user := auth.UserFromContext(ctx)

		pokemonID := c.Param("pokemonId")
		if pokemonID == "" {
			log.Debug().Msg("Pokemon ID is missing from path")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Pokemon ID is missing"})
		}

		pokemon, err := models.CollectionPokemons(
			models.CollectionPokemonWhere.PokemonID.EQ(pokemonID),
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
		).One(ctx, s.DB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to retrieve pokemon from collection")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}

		if pokemon == nil {
			log.Debug().Str("pokemonID", pokemonID).Msg("Pokemon does not exist in collection")
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Pokemon not found in collection"})
		}

		var updateData struct {
			Caught   *time.Time `json:"caught"`
			Nickname *string    `json:"nickname"`
		}
		if err := c.Bind(&updateData); err != nil {
			log.Debug().Err(err).Msg("Invalid request body")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if updateData.Caught != nil {
			pokemon.Caught = null.TimeFrom(*updateData.Caught)
		}
		if updateData.Nickname != nil {
			pokemon.Nickname = null.StringFrom(*updateData.Nickname)
		}

		_, err = pokemon.Update(ctx, s.DB, boil.Whitelist("caught", "nickname"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to update pokemon in collection")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update pokemon"})
		}

		log.Info().Str("pokemonID", pokemonID).Str("userID", user.ID).Msg("Pokemon successfully updated")
		return c.JSON(http.StatusOK, map[string]string{"message": "Pokemon successfully updated"})
	}
}
