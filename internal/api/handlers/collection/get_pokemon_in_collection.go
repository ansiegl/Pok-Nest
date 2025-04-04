package collection

import (
	"net/http"
	"strings"

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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

		if params.Name != nil {
			nameFormatted := cases.Title(language.English).String(strings.ToLower(*params.Name))
			queryMods = append(queryMods, qm.Where("name = ?", nameFormatted))
		}
		if params.Type != nil {
			typeFormatted := cases.Title(language.English).String(strings.ToLower(*params.Type))
			queryMods = append(queryMods, qm.Where("type_1 = ? OR type_2 = ?", typeFormatted, typeFormatted))
		}
		if params.Hp != nil {
			queryMods = append(queryMods, qm.Where("hp = ?", params.Hp))
		}
		if params.Attack != nil {
			queryMods = append(queryMods, qm.Where("attack = ?", params.Attack))
		}
		if params.Hp != nil {
			queryMods = append(queryMods, qm.Where("hp = ?", params.Hp))
		}
		if params.Defense != nil {
			queryMods = append(queryMods, qm.Where("defense = ?", params.Defense))
		}
		if params.Speed != nil {
			queryMods = append(queryMods, qm.Where("speed = ?", params.Speed))
		}
		if params.Special != nil {
			queryMods = append(queryMods, qm.Where("special = ?", params.Special))
		}

		sortOrder := "asc"
		if params.SortOrder != nil {
			sortOrder = *params.SortOrder
		}
		queryMods = append(queryMods, qm.OrderBy("pokemon.name "+sortOrder))

		limit := 10
		if params.Limit != nil {
			limit = int(*params.Limit)
		}
		offset := 0
		if params.Offset != nil {
			offset = int(*params.Offset)
		}

		queryMods = append(queryMods, qm.Limit(limit), qm.Offset(offset))

		totalCount, err := models.CollectionPokemons().Count(ctx, s.DB)
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

				var nameOrNickname *string
				if collectionPokemon.Nickname.Valid && collectionPokemon.Nickname.String != "" {
					nameOrNickname = &collectionPokemon.Nickname.String
				} else {
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
					Png:            &pokemon.PNGURL,
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
