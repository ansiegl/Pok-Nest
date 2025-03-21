package pokemon

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/labstack/echo/v4"
)

func GetPokemonDownloadRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Pokemon.GET("/download", getPokemonDownloadHandler(s))
}

func getPokemonDownloadHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		pokemons, err := models.Pokemons().All(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Failed to load Pokémon")
			return err
		}

		csvData := [][]string{
			{"PokemonID", "Name", "Type1", "Type2", "Generation", "Legendary"},
		}

		for _, p := range pokemons {
			csvData = append(csvData, []string{
				p.PokemonID,
				p.Name,
				p.Type1,
				p.Type2.String,
				strconv.Itoa(p.Generation),
				boolToString(p.Legendary),
			})
		}

		c.Response().Header().Set("Content-Type", "text/csv")
		c.Response().Header().Set("Content-Disposition", "attachment; filename=pokemon.csv")

		w := csv.NewWriter(c.Response().Writer)
		err = w.WriteAll(csvData)
		if err != nil {
			log.Err(err).Msg("Failed to write CSV")
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "CSV download successful"})
	}
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
