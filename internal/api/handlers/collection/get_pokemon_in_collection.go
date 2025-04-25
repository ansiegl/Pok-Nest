package collection

import (
	"net/http"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/auth"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/ansiegl/Pok-Nest.git/internal/types/collection"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetCollectionPokemonRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Collection.GET("/pokemon", getCollectionPokemonHandler(s))
}

func getCollectionPokemonHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		user := auth.UserFromContext(ctx)

		// get parameter
		params := collection.NewGetCollectionPokemonParams()
		if err := util.BindAndValidatePathAndQueryParams(c, &params); err != nil {
			return err
		}

		// create query
		queryMods := []qm.QueryMod{
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
			qm.InnerJoin("pokemon ON pokemon.pokemon_id = collection_pokemon.pokemon_id"),
			qm.Load("Pokemon"),
		}

		limit := 10
		if params.Limit != nil {
			limit = int(*params.Limit)
		}
		offset := 0
		if params.Offset != nil {
			offset = int(*params.Offset)
		}

		queryMods = append(queryMods, qm.Limit(limit), qm.Offset(offset))

		totalCount, err := models.CollectionPokemons(qm.Where("user_id = ?", user.ID)).Count(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Failed to get total count of collection pokemon")
			return err
		}

		collectionPokemons, err := models.CollectionPokemons(queryMods...).All(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Failed to load pokemon")
			return err
		}

		var pokemonData []*types.CollectionPokemon
		for _, collectionPokemon := range collectionPokemons {
			pokemon := collectionPokemon.R.Pokemon
			if pokemon != nil {
				pokemonID := strfmt.UUID4(pokemon.PokemonID)

				nameOrNickname := util.NonEmptyOrNil(collectionPokemon.Nickname.String)
				if nameOrNickname == nil {
					nameOrNickname = &pokemon.Name
				}

				pokemonData = append(pokemonData, &types.CollectionPokemon{
					PokemonID:      &pokemonID,
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
				})
			}
		}

		tempLimit := int64(limit)
		tempOffset := int64(offset)

		response := &types.GetCollectionPokemonResponse{
			Data: pokemonData,
			Pagination: &types.Pagination{
				Total:  totalCount,
				Limit:  tempLimit,
				Offset: tempOffset,
			},
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
