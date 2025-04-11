package auth

import (
	"database/sql"
	"errors"
	"net/http"

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

func GetUserInfoRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.GET("/userinfo", getUserInfoHandler(s))
}

func getUserInfoHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := auth.UserFromContext(ctx)
		log := util.LogFromContext(ctx)

		totalCollected, err := models.CollectionPokemons(
			qm.Where("user_id = ?", user.ID),
		).Count(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("failed to count collected pokemon")
			return err
		}

		totalInDB, err := models.Pokemons().Count(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("failed to count all pokemon")
			return err
		}
		missingFromDatabase := totalInDB - totalCollected

		// separate maps for first and sec type
		primaryTypeStats := make(map[string]int64)
		secondaryTypeStats := make(map[string]int64)

		type typeCount struct {
			Type  string `boil:"type_1"`
			Count int64  `boil:"count"`
		}
		var primaryTypeCounts []typeCount

		err = models.NewQuery(
			qm.Select("p.type_1, COUNT(*) as count"),
			qm.From("collection_pokemon cp"),
			qm.InnerJoin("pokemon p ON cp.pokemon_id = p.pokemon_id"),
			qm.Where("cp.user_id = ?", user.ID),
			qm.GroupBy("p.type_1"),
		).Bind(ctx, s.DB, &primaryTypeCounts)
		if err != nil {
			log.Err(err).Msg("failed to get primary type statistics")
			return err
		}

		for _, tc := range primaryTypeCounts {
			primaryTypeStats[tc.Type] = tc.Count
		}

		type secTypeCount struct {
			Type  string `boil:"type_2"`
			Count int64  `boil:"count"`
		}
		var secondaryTypeCounts []secTypeCount

		err = models.NewQuery(
			qm.Select("p.type_2, COUNT(*) as count"),
			qm.From("collection_pokemon cp"),
			qm.InnerJoin("pokemon p ON cp.pokemon_id = p.pokemon_id"),
			qm.Where("cp.user_id = ? AND p.type_2 IS NOT NULL", user.ID),
			qm.GroupBy("p.type_2"),
		).Bind(ctx, s.DB, &secondaryTypeCounts)
		if err != nil {
			log.Err(err).Msg("failed to get secondary type statistics")
			return err
		}

		for _, tc := range secondaryTypeCounts {
			secondaryTypeStats[tc.Type] = tc.Count
		}

		response := &types.GetUserInfoResponse{
			Sub:       swag.String(user.ID),
			UpdatedAt: swag.Int64(user.UpdatedAt.Unix()),
			Email:     strfmt.Email(user.Username.String),
			Scopes:    user.Scopes,
			CollectionStats: &types.GetUserInfoResponseCollectionStats{
				TotalCollected:      *swag.Int64(totalCollected),
				MissingFromDatabase: *swag.Int64(missingFromDatabase),
				FirstTypes:          primaryTypeStats,
				SecTypes:            secondaryTypeStats,
			},
		}

		// if this user has an appUserProfile attached, add additional / modify props from there
		appUserProfile, err := user.AppUserProfile().One(ctx, s.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return util.ValidateAndReturn(c, http.StatusOK, response)
			}

			log.Debug().Err(err).Msg("Unknown error while getting appUserProfile information for user")
			return err
		}

		if appUserProfile.UpdatedAt.After(user.UpdatedAt) {
			response.UpdatedAt = swag.Int64(appUserProfile.UpdatedAt.Unix())
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
