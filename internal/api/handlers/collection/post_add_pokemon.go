package collection

import (
	"net/http"
	"time"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/auth"
	"github.com/ansiegl/Pok-Nest.git/internal/api/httperrors"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/ansiegl/Pok-Nest.git/internal/types/collection"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/ansiegl/Pok-Nest.git/internal/util/db"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func PostAddPokemonToCollectionRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Pokemon.POST("/:pokemonId", postAddPokemonToCollectionHandler(s))
}

func postAddPokemonToCollectionHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)
		user := auth.UserFromContext(ctx)

		params := collection.NewPostAddPokemonToCollectionParams()
		err := util.BindAndValidatePathAndQueryParams(c, &params)
		if err != nil {
			return err
		}

		pokemonExists, err := models.Pokemons(models.PokemonWhere.PokemonID.EQ(params.PokemonID)).Exists(ctx, s.DB)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to check if pokemon exists")
			return err
		}
		if !pokemonExists {
			log.Debug().Str("pokemonID", params.PokemonID).Msg("Pokemon does not exist")
			return httperrors.NewHTTPError(http.StatusNotFound, "POKEMON_NOT_FOUND", "Pokemon not found")
		}

		var request types.PokemonBody
		if err := c.Bind(&request); err != nil {
			log.Debug().Err(err).Msg("Failed to bind request")
			return httperrors.NewHTTPError(http.StatusBadRequest, "INVALID_BODY", "Invalid request body")
		}

		caughtTime, err := time.Parse("2006-01-02", request.Caught.String())
		if err != nil {
			log.Debug().Err(err).Str("caught", request.Caught.String()).Msg("Failed to parse caught date")
			return httperrors.NewHTTPError(http.StatusBadRequest, "INVALID_DATE", "Invalid date format")
		}

		if err := db.WithTransaction(ctx, s.DB, func(tx boil.ContextExecutor) error {
			collection := &models.CollectionPokemon{
				UserID:    user.ID,
				PokemonID: params.PokemonID,
				Caught:    null.TimeFrom(caughtTime),
				Nickname:  null.StringFrom(request.Nickname),
			}
			if err := collection.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert collection")
				return err
			}
			log.Debug().Str("pokemonID", params.PokemonID).Str("userID", user.ID).Msg("Successfully added pokemon to collection")
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
