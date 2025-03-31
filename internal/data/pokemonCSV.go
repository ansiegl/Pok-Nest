package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/volatiletech/null/v8"
)

func LoadPokemonFromCSV(csvPath string) ([]models.Pokemon, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var pokemons []models.Pokemon
	lineNumber := 0

	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("Error reading CSV header: %w", err)
	}

	for {
		lineNumber++
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("Error in line %d: %v\n", lineNumber, err)
			continue
		}

		var type2 null.String
		if row[3] != "" {
			type2 = null.StringFrom(row[3])
		} else {
			type2 = null.String{Valid: false}
		}

		gen, err := strconv.Atoi(row[11])
		if err != nil {
			fmt.Printf("Error in line %d: Invalid 'Generation' value (%v)\n", lineNumber, row[11])
			continue
		}

		legendary, err := strconv.ParseBool(row[12])
		if err != nil {
			fmt.Printf("Error in line %d: Invalid 'Legendary' value (%v)\n", lineNumber, row[12])
			continue
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
