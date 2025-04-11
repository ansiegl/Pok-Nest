package pokemon

import (
	"net/http"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/types"
	"github.com/ansiegl/Pok-Nest.git/internal/types/pokemon"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetPokemonRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Pokemon.GET("", getPokemonHandler(s))
}

func getPokemonHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		params := pokemon.NewGetAllPokemonParams()
		err := util.BindAndValidatePathAndQueryParams(c, &params)
		if err != nil {
			return err
		}

		queryMods := []qm.QueryMod{}

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
