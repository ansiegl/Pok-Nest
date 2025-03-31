package pokemon

import (
	"bytes"
	"encoding/csv"
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

		var csvBuffer bytes.Buffer
		writer := csv.NewWriter(&csvBuffer)

		header := []string{"PokemonID", "Name", "Type1", "Type2", "Generation", "Legendary"}
		if err := writer.Write(header); err != nil {
			log.Err(err).Msg("Failed to write CSV header")
			return err
		}

		for _, p := range pokemons {
			row := []string{
				p.PokemonID,
				p.Name,
				p.Type1,
				p.Type2.String,
				strconv.Itoa(p.Generation),
				boolToString(p.Legendary),
			}
			if err := writer.Write(row); err != nil {
				log.Err(err).Msg("Failed to write CSV row")
				return err
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Err(err).Msg("Failed to flush CSV writer")
			return err
		}

		c.Response().Header().Set("Content-Type", "text/csv")
		c.Response().Header().Set("Content-Disposition", "attachment; filename=pokemon.csv")

		_, err = c.Response().Write(csvBuffer.Bytes())
		if err != nil {
			log.Err(err).Msg("Failed to write response")
			return err
		}

		return nil
	}
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
