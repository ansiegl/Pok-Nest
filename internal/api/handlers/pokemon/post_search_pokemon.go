package pokemon

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/ansiegl/Pok-Nest.git/internal/types/pokemon"
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

		params := pokemon.NewPostSearchPokemonParams()
		err := util.BindAndValidatePathAndQueryParams(c, &params)
		if err != nil {
			return err
		}

		var searchRequest types.PokemonSearchRequest
		if err := c.Bind(&searchRequest); err != nil {
			log.Err(err).Msg("Failed to bind request body")
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		queryMods := []qm.QueryMod{}

		if searchRequest.Name != "" {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s LIKE ?", models.PokemonTableColumns.Name), "%"+searchRequest.Name+"%"),
			)
		}
		if searchRequest.Type != "" {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("(%s = ? OR %s = ?)",
					models.PokemonTableColumns.Type1, models.PokemonTableColumns.Type2,
				), searchRequest.Type, searchRequest.Type),
			)
		}
		if searchRequest.Hp != nil && *searchRequest.Hp > 0 {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s >= ?", models.PokemonTableColumns.HP), *searchRequest.Hp),
			)
		}
		if searchRequest.Attack != nil && *searchRequest.Attack > 0 {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s >= ?", models.PokemonTableColumns.Attack), *searchRequest.Attack),
			)
		}
		if searchRequest.Defense != nil && *searchRequest.Defense > 0 {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s >= ?", models.PokemonTableColumns.Defense), *searchRequest.Defense),
			)
		}
		if searchRequest.Speed != nil && *searchRequest.Speed > 0 {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s >= ?", models.PokemonTableColumns.Speed), *searchRequest.Speed),
			)
		}
		if searchRequest.Special != nil && *searchRequest.Special > 0 {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s >= ?", models.PokemonTableColumns.Special), *searchRequest.Special),
			)
		}
		if searchRequest.SortOrder != "" {
			orderCol := models.PokemonTableColumns.Name
			if strings.ToLower(searchRequest.SortOrder) == "asc" {
				queryMods = append(queryMods, qm.OrderBy(orderCol+" ASC"))
			} else if strings.ToLower(searchRequest.SortOrder) == "desc" {
				queryMods = append(queryMods, qm.OrderBy(orderCol+" DESC"))
			}
		}

		totalCount, err := models.Pokemons(queryMods...).Count(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Failed to get total count of pokemon")
			return err
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

		response := &types.GetPokemonResponse{
			Data: pokemonData,
			Pagination: &types.Pagination{
				Total:  totalCount,
				Limit:  int64(limit),
				Offset: int64(offset),
			},
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
