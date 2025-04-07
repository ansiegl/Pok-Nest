package collection

import (
	"net/http"
	"strings"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/auth"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func PostSearchPokemonInCollectionRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Collection.POST("/pokemon", postSearchPokemonInCollectionHandler(s))
}

func postSearchPokemonInCollectionHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		user := auth.UserFromContext(ctx)

		var searchRequest types.PokemonSearchRequest
		if err := c.Bind(&searchRequest); err != nil {
			log.Err(err).Msg("Failed to bind request body")
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		queryMods := []qm.QueryMod{
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
			qm.InnerJoin("pokemon ON pokemon.pokemon_id = collection_pokemon.pokemon_id"),
			qm.Load("Pokemon"),
		}

		if searchRequest.Name != "" {
			queryMods = append(queryMods, qm.Where("name LIKE ?", "%"+searchRequest.Name+"%"))
		}
		if searchRequest.Type != "" {
			queryMods = append(queryMods, qm.Where("type_1 = ? OR type_2 = ?", searchRequest.Type, searchRequest.Type))
		}
		if searchRequest.Hp > 0 {
			queryMods = append(queryMods, qm.Where("hp >= ?", searchRequest.Hp))
		}
		if searchRequest.Attack > 0 {
			queryMods = append(queryMods, qm.Where("attack >= ?", searchRequest.Attack))
		}
		if searchRequest.Defense > 0 {
			queryMods = append(queryMods, qm.Where("defense >= ?", searchRequest.Defense))
		}
		if searchRequest.Speed > 0 {
			queryMods = append(queryMods, qm.Where("speed >= ?", searchRequest.Speed))
		}
		if searchRequest.Special > 0 {
			queryMods = append(queryMods, qm.Where("special >= ?", searchRequest.Special))
		}
		if searchRequest.SortOrder != "" {
			if strings.ToLower(searchRequest.SortOrder) == "asc" {
				queryMods = append(queryMods, qm.OrderBy("name ASC"))
			} else if strings.ToLower(searchRequest.SortOrder) == "desc" {
				queryMods = append(queryMods, qm.OrderBy("name DESC"))
			}
		}

		limit := 10
		offset := 0

		if searchRequest.Pagination != nil {
			if searchRequest.Pagination.Limit > 0 {
				limit = int(searchRequest.Pagination.Limit)
			}
			if searchRequest.Pagination.Offset > 0 {
				offset = int(searchRequest.Pagination.Offset)
			}
		}

		queryMods = append(queryMods, qm.Limit(limit), qm.Offset(offset))

		totalCount, err := models.CollectionPokemons(queryMods...).Count(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Failed to get total count of pokemon")
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

				pokemonData = append(pokemonData, &types.CollectionPokemon{
					PokemonID:      &pokemonID,
					Number:         swag.Int64(int64(pokemon.PokemonNumber)),
					NameOrNickname: &pokemon.Name,
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
