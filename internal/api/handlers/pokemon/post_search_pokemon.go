package pokemon

import (
	"net/http"
	"strings"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func PostSearchPokemonRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Pokemon.POST("", postSearchPokemonHandler(s))
}

func postSearchPokemonHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var searchRequest types.PokemonSearchRequest
		if err := c.Bind(&searchRequest); err != nil {
			log.Err(err).Msg("Failed to bind request body")
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		queryMods := []qm.QueryMod{}

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

		totalCount, err := models.Pokemons(queryMods...).Count(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Failed to get total count of pokemon")
			return err
		}

		pokemons, err := models.Pokemons(queryMods...).All(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Failed to load pokemon")
			return err
		}

		var pokemonData []*types.Pokemon
		for _, p := range pokemons {
			pokemonID := strfmt.UUID4(p.PokemonID)

			pokemonData = append(pokemonData, &types.Pokemon{
				PokemonID:   &pokemonID,
				Number:      int64(p.PokemonNumber),
				Name:        &p.Name,
				Type1:       &p.Type1,
				Type2:       p.Type2.String,
				Hp:          int64(p.HP),
				Attack:      int64(p.Attack),
				Defense:     int64(p.Defense),
				Speed:       int64(p.Speed),
				Special:     int64(p.Special),
				ImageURL:    p.PNGURL,
				Description: p.Description,
			})
		}

		tempLimit := int64(limit)
		tempOffset := int64(offset)

		response := &types.GetPokemonResponse{
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
