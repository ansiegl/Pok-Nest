package pokemon

import (
	"net/http"
	"strconv"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/labstack/echo/v4"
)

func GetPokemonRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Pokemon.GET("", getPokemonHandler(s))
}

func getPokemonHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		// get limit from request
		limitStr := c.QueryParam("limit")
		limit := 10 // default value

		if limitStr != "" {
			parsedLimit, err := strconv.Atoi(limitStr)
			if err != nil {
				log.Err(err).Str("limit", limitStr).Msg("Invalid limit parameter")
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid limit parameter",
				})
			}

			if parsedLimit <= 0 {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Limit must be greater than 0",
				})
			}

			if parsedLimit > 20 {
				// set maximum
				limit = 20
			} else {
				limit = parsedLimit
			}
		}

		// get offset from request
		offsetStr := c.QueryParam("offset")
		offset := 0 // default value

		if offsetStr != "" {
			parsedOffset, err := strconv.Atoi(offsetStr)
			if err != nil {
				log.Err(err).Str("offset", offsetStr).Msg("Invalid offset parameter")
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid offset parameter",
				})
			}

			if parsedOffset < 0 {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Offset must be greater than or equal to 0"})
			}

			offset = parsedOffset
		}

		pokemons, err := models.Pokemons().All(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Failed to load pokemon")
			return err
		}

		if offset >= len(pokemons) {
			return c.JSON(http.StatusOK, []struct{}{})
		}

		// check number of pokemon to return
		endIndex := offset + limit
		if endIndex > len(pokemons) {
			endIndex = len(pokemons)
		}

		// specific order in response
		type PokemonResponse struct {
			PokemonID  string `json:"pokemon_id"`
			Name       string `json:"name"`
			Type1      string `json:"type_1"`
			Type2      string `json:"type_2,omitempty"`
			Generation int    `json:"generation"`
			Legendary  bool   `json:"legendary"`
		}

		// pagination metadata
		type PaginationMetadata struct {
			Total  int `json:"total"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}

		// response with pagination metadata
		type APIResponse struct {
			Data       []PokemonResponse  `json:"data"`
			Pagination PaginationMetadata `json:"pagination"`
		}

		var responseData []PokemonResponse
		for _, p := range pokemons[offset:endIndex] {
			responseData = append(responseData, PokemonResponse{
				PokemonID:  p.PokemonID,
				Name:       p.Name,
				Type1:      p.Type1,
				Type2:      p.Type2.String,
				Generation: p.Generation,
				Legendary:  p.Legendary,
			})
		}

		response := APIResponse{
			Data: responseData,
			Pagination: PaginationMetadata{
				Total:  len(pokemons),
				Limit:  limit,
				Offset: offset,
			},
		}

		return c.JSON(http.StatusOK, response)
	}
}
