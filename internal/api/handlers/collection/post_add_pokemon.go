package collection

import (
	"net/http"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/auth"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/ansiegl/Pok-Nest.git/internal/util/db"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func PostAddPokemonToCollectionRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Pokemon.POST("/:pokemonId", postAddPokemonToCollectionHandler(s))
}

func postAddPokemonToCollectionHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := auth.UserFromContext(ctx)
		log := util.LogFromContext(ctx)

		pokemonID := c.Param("pokemonId")
		if pokemonID == "" {
			log.Debug().Msg("Pokemon ID is missing from path")
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Pokemon ID is missing"})
		}

		pokemonExists, err := models.Pokemons(models.PokemonWhere.PokemonID.EQ(pokemonID)).Exists(ctx, s.DB)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to check if pokemon exsists")
		}

		if !pokemonExists {
			log.Debug().Str("pokemonID", pokemonID).Msg("Pokemon does not exist")
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Pokemon not found"})
		}

		var request models.CollectionPokemon
		if err := c.Bind(&request); err != nil {
			log.Debug().Err(err).Msg("Failed to bind request")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if err := db.WithTransaction(ctx, s.DB, func(tx boil.ContextExecutor) error {
			collection := &models.CollectionPokemon{
				UserID:    user.ID,
				PokemonID: pokemonID,
				Caught:    request.Caught,
				Nickname:  request.Nickname,
			}

			if err := collection.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert collection")
				return err
			}

			log.Debug().Str("pokemonID", pokemonID).Str("userID", user.ID).Msg("Successfully added pokemon to collection")
			return nil
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to add pokemon to collection")
			return err
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"message": "Pokemon successfully added to collection",
		})

	}
}
