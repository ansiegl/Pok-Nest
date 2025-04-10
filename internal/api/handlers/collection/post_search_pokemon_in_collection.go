package collection

import (
	"fmt"
	"net/http"
	"strings"

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

func PostSearchPokemonInCollectionRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Collection.POST("/pokemon", postSearchPokemonInCollectionHandler(s))
}

func postSearchPokemonInCollectionHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)
		user := auth.UserFromContext(ctx)

		params := collection.NewPostSearchPokemonInCollectionParams()
		err := util.BindAndValidatePathAndQueryParams(c, &params)
		if err != nil {
			return err
		}

		var searchRequest types.PokemonSearchRequest
		if err := c.Bind(&searchRequest); err != nil {
			log.Err(err).Msg("Failed to bind request body")
			return httperrors.NewHTTPError(http.StatusBadRequest, "INVALID_BODY", "Invalid request body")
		}

		queryMods := []qm.QueryMod{
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
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

		if searchRequest.Name != "" {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s LIKE ?", models.PokemonTableColumns.Name), "%"+searchRequest.Name+"%"),
			)
		}
		if searchRequest.Type != "" {
			queryMods = append(queryMods,
				qm.Where(
					fmt.Sprintf("(%s = ? OR %s = ?)", models.PokemonTableColumns.Type1, models.PokemonTableColumns.Type2),
					searchRequest.Type, searchRequest.Type,
				),
			)
		}
		if searchRequest.Hp != nil && *searchRequest.Hp > 0 {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s >= ?", models.PokemonTableColumns.HP), searchRequest.Hp),
			)
		}
		if searchRequest.Attack != nil && *searchRequest.Attack > 0 {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s >= ?", models.PokemonTableColumns.Attack), searchRequest.Attack),
			)
		}
		if searchRequest.Defense != nil && *searchRequest.Defense > 0 {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s >= ?", models.PokemonTableColumns.Defense), searchRequest.Defense),
			)
		}
		if searchRequest.Speed != nil && *searchRequest.Speed > 0 {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s >= ?", models.PokemonTableColumns.Speed), searchRequest.Speed),
			)
		}
		if searchRequest.Special != nil && *searchRequest.Special > 0 {
			queryMods = append(queryMods,
				qm.Where(fmt.Sprintf("%s >= ?", models.PokemonTableColumns.Special), searchRequest.Special),
			)
		}
		if searchRequest.SortOrder != "" {
			orderColumn := models.PokemonTableColumns.Name
			if strings.ToLower(searchRequest.SortOrder) == "asc" {
				queryMods = append(queryMods, qm.OrderBy(orderColumn+" ASC"))
			} else if strings.ToLower(searchRequest.SortOrder) == "desc" {
				queryMods = append(queryMods, qm.OrderBy(orderColumn+" DESC"))
			}
		}

		totalCount, err := models.CollectionPokemons(queryMods...).Count(ctx, s.DB)
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

		collectionPokemons, err := models.CollectionPokemons(queryMods...).All(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Failed to load pokemon")
			return err
		}

		var pokemonData []*types.CollectionPokemon
		for _, cp := range collectionPokemons {
			p := cp.R.Pokemon
			if p == nil {
				continue
			}
			pokemonID := strfmt.UUID4(p.PokemonID)
			pokemonData = append(pokemonData, &types.CollectionPokemon{
				PokemonID:      &pokemonID,
				Number:         swag.Int64(int64(p.PokemonNumber)),
				NameOrNickname: &p.Name,
				Type1:          &p.Type1,
				Type2:          p.Type2.String,
				Hp:             swag.Int64(int64(p.HP)),
				Attack:         swag.Int64(int64(p.Attack)),
				Defense:        swag.Int64(int64(p.Defense)),
				Speed:          swag.Int64(int64(p.Speed)),
				Special:        swag.Int64(int64(p.Special)),
				ImageURL:       &p.PNGURL,
				Description:    &p.Description,
			})
		}

		response := &types.GetCollectionPokemonResponse{
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
