package data

import (
	"fmt"

	"encoding/csv" // read csv files
	"os"           // open and close files
	"strconv"      // convert datatypes

	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/volatiletech/null/v8"
)

// read csv
func LoadPokemonFromCSV() ([]models.Pokemon, error) {
	file, err := os.Open("docs/pokemon.csv") // open csv file
	if err != nil {
		return nil, fmt.Errorf("Error opening CSV file: %w", err)
	}
	defer file.Close() // close csv file

	reader := csv.NewReader(file)
	records, err := reader.ReadAll() // reads csv file and saves it as string lists

	if err != nil {
		return nil, fmt.Errorf("Error while reading CSV file: %w", err)
	}

	var pokemons []models.Pokemon

	for i, row := range records {

		if i == 0 { // skip first row
			continue
		}

		var type2 null.String // type2 can be null
		if row[3] != "" {
			type2 = null.StringFrom(row[3])
		} else {
			type2 = null.String{Valid: false}
		}

		gen, err := strconv.Atoi(row[11])
		if err != nil {
			return nil, fmt.Errorf("Error converting 'Generation' into an integer: %w", err)
		}

		legendary, err := strconv.ParseBool(row[12])
		if err != nil {
			return nil, fmt.Errorf("Error parsing 'Legendary' value: %w", err)
		}

		pokemon := models.Pokemon{
			Name:       row[1],
			Type1:      row[2],
			Type2:      type2,
			Generation: gen,
			Legendary:  legendary,
		}

		pokemons = append(pokemons, pokemon)
	}
	return pokemons, nil
}
