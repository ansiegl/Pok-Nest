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
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetCollectionPokemonDetailRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Collection.GET("/pokemon/:pokemonId", getCollectionPokemonDetailHandler(s))
}

func getCollectionPokemonDetailHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)
		user := auth.UserFromContext(ctx)

		params := collection.NewGetCollectionPokemonDetailParams()
		if err := util.BindAndValidatePathAndQueryParams(c, &params); err != nil {
			return err
		}

		queryMods := []qm.QueryMod{
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
			models.CollectionPokemonWhere.PokemonID.EQ(params.PokemonID),

			qm.InnerJoin(
				fmt.Sprintf(
					"%s ON %s = %s",
					models.TableNames.Pokemon,
					models.PokemonTableColumns.PokemonID,
					models.CollectionPokemonTableColumns.PokemonID,
				),
			),
			qm.Load(qm.Rels(models.CollectionPokemonRels.Pokemon)),
		}

		collectionPokemon, err := models.CollectionPokemons(queryMods...).One(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Pokemon not found")
			return httperrors.NewHTTPError(http.StatusNotFound, "POKEMON_NOT_FOUND", fmt.Sprintf("No pokemon found that matches the ID %v", params.PokemonID))
		}

		pokemon := collectionPokemon.R.Pokemon
		if pokemon == nil {
			return httperrors.NewHTTPError(http.StatusNotFound, "POKEMONDETAILS_NOT_FOUND", fmt.Sprintf("Pokemon details for %v not available", params.PokemonID))
		}

		var caughtDate strfmt.Date
		if collectionPokemon.Caught.Valid {
			caughtDate = strfmt.Date(collectionPokemon.Caught.Time)
		}

		var nameOrNickname *string
		if collectionPokemon.Nickname.Valid && collectionPokemon.Nickname.String != "" {
			nameOrNickname = &collectionPokemon.Nickname.String
		} else {
			nameOrNickname = &pokemon.Name
		}

		pokemonIDStr := strfmt.UUID4(pokemon.PokemonID)
		response := &types.CollectionPokemonDetail{
			CollectionPokemon: types.CollectionPokemon{
				PokemonID:      &pokemonIDStr,
				Number:         swag.Int64(int64(pokemon.PokemonNumber)),
				NameOrNickname: nameOrNickname,
				Type1:          &pokemon.Type1,
				Type2:          pokemon.Type2.String,
				Hp:             swag.Int64(int64(pokemon.HP)),
				Attack:         swag.Int64(int64(pokemon.Attack)),
				Defense:        swag.Int64(int64(pokemon.Defense)),
				Speed:          swag.Int64(int64(pokemon.Speed)),
				Special:        swag.Int64(int64(pokemon.Special)),
				ImageURL:       &pokemon.PNGURL,
				Description:    &pokemon.Description,
			},
			GifURL: pokemon.GifURL,
			Caught: &caughtDate,
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
