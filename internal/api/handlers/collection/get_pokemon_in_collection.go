package collection

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/auth"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/res"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetPokemonRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Collection.GET("/pokemon", getPokemonHandler(s))
}

func getPokemonHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		user := auth.UserFromContext(ctx)

		// get name from request
		name := c.QueryParam("name")

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

		// get filter parameters from request
		pokemonType := c.QueryParam("type")
		generationStr := c.QueryParam("generation")
		legendaryStr := c.QueryParam("legendary")

		var generation int
		if generationStr != "" {
			var err error
			generation, err = strconv.Atoi(generationStr)
			if err != nil {
				log.Err(err).Str("generation", generationStr).Msg("Invalid generation parameter")
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid generation parameter",
				})
			}
		}

		legendary := false
		if legendaryStr != "" {
			legendary = legendaryStr == "true"
		}

		// get sortOrder from request (asc or desc)
		sortOrder := c.QueryParam("sortOrder")
		if sortOrder != "asc" && sortOrder != "desc" && sortOrder != "" {
			log.Debug().Str("sortOrder", sortOrder).Msg("Invalid sortOrder parameter")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sortOrder parameter. Allowed values are 'asc' or 'desc'."})
		}

		collectionPokemons, err := models.CollectionPokemons(
			models.CollectionPokemonWhere.UserID.EQ(user.ID),
			qm.Load(models.CollectionPokemonRels.Pokemon),
		).All(ctx, s.DB)
		if err != nil {
			log.Err(err).Msg("Fehler beim Laden der Pokémon")
			return err
		}

		// filter pokemon based on query parameters
		var filteredPokemons []models.Pokemon
		for _, collectionPokemon := range collectionPokemons {
			if collectionPokemon.R != nil && collectionPokemon.R.Pokemon != nil {
				p := collectionPokemon.R.Pokemon
				// filter by name
				if name != "" && p.Name != name {
					continue
				}

				// filter by type (either type 1 or type 2)
				if pokemonType != "" && (p.Type1 != pokemonType && p.Type2.String != pokemonType) {
					continue
				}

				// filter by generation
				if generation > 0 && p.Generation != generation {
					continue
				}

				// filter by legendary
				if legendary && !p.Legendary {
					continue
				}

				// add pokemon to filtered list
				filteredPokemons = append(filteredPokemons, *p)
			}
		}

		// Sort the filteredPokemons based on sortOrder
		if sortOrder == "asc" {
			// Ascending sort by Name (you can change to another field like Generation if needed)
			sort.Slice(filteredPokemons, func(i, j int) bool {
				return filteredPokemons[i].Name < filteredPokemons[j].Name
			})
		} else if sortOrder == "desc" {
			// Descending sort by Name (you can change to another field like Generation if needed)
			sort.Slice(filteredPokemons, func(i, j int) bool {
				return filteredPokemons[i].Name > filteredPokemons[j].Name
			})
		}

		// pagination logic
		if offset >= len(filteredPokemons) {
			return c.JSON(http.StatusOK, []struct{}{})
		}

		endIndex := offset + limit
		if endIndex > len(filteredPokemons) {
			endIndex = len(filteredPokemons)
		}

		var responseData []res.PokemonResponse
		for _, collectionPokemon := range collectionPokemons[offset:endIndex] {
			if collectionPokemon.R != nil && collectionPokemon.R.Pokemon != nil {
				pokemon := collectionPokemon.R.Pokemon
				responseData = append(responseData, res.PokemonResponse{
					PokemonID:  pokemon.PokemonID,
					Name:       pokemon.Name,
					Type1:      pokemon.Type1,
					Type2:      pokemon.Type2.String,
					Generation: pokemon.Generation,
					Legendary:  pokemon.Legendary,
				})
			}
		}

		// return response with pagination metadata
		response := res.APIResponse{
			Data: responseData,
			Pagination: res.PaginationMetadata{
				Total:  len(filteredPokemons),
				Limit:  limit,
				Offset: offset,
			},
		}

		return c.JSON(http.StatusOK, response)
	}
}
