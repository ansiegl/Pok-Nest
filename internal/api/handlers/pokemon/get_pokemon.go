package pokemon

import (
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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GetPokemonRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Pokemon.GET("", getPokemonHandler(s))
}

func getPokemonHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		params := pokemon.NewGetPokemonAsJSONParams()
		err := util.BindAndValidatePathAndQueryParams(c, &params)
		if err != nil {
			return err
		}

		queryMods := []qm.QueryMod{}

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
		queryMods = append(queryMods, qm.OrderBy("name "+sortOrder))

		limit := 10
		if params.Limit != nil {
			limit = int(*params.Limit)
		}
		offset := 0
		if params.Offset != nil {
			offset = int(*params.Offset)
		}

		queryMods = append(queryMods, qm.Limit(limit), qm.Offset(offset))

		totalCount, err := models.Pokemons().Count(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Failed to get total count of Pokémon")
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
				Name:        &p.Name,
				Type1:       &p.Type1,
				Type2:       p.Type2.String,
				Hp:          int64(p.HP),
				Attack:      int64(p.Attack),
				Defense:     int64(p.Defense),
				Speed:       int64(p.Speed),
				Special:     int64(p.Special),
				Png:         p.PNGURL,
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
