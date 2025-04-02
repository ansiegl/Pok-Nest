package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/volatiletech/null/v8"
)

func validateStats(statStr string) (int, error) {
	stat, err := strconv.Atoi(statStr)
	if err != nil || stat < 0 {
		return 0, fmt.Errorf("invalid stat value: %v", statStr)
	}
	return stat, nil
}

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

		hp, err := validateStats(row[4])
		if err != nil {
			fmt.Printf("Error in line %d: Invalid 'HP' value (%v)\n", lineNumber, row[4])
			continue
		}

		attack, err := validateStats(row[5])
		if err != nil {
			fmt.Printf("Error in line %d: Invalid 'Attack' value (%v)\n", lineNumber, row[5])
			continue
		}

		defense, err := validateStats(row[6])
		if err != nil {
			fmt.Printf("Error in line %d: Invalid 'Defense' value (%v)\n", lineNumber, row[6])
			continue
		}

		speed, err := validateStats(row[7])
		if err != nil {
			fmt.Printf("Error in line %d: Invalid 'Speed' value (%v)\n", lineNumber, row[7])
			continue
		}

		special, err := validateStats(row[8])
		if err != nil {
			fmt.Printf("Error in line %d: Invalid 'Special' value (%v)\n", lineNumber, row[8])
			continue
		}

		pokemon := models.Pokemon{
			Name:        row[1],
			Type1:       row[2],
			Type2:       type2,
			HP:          hp,
			Attack:      attack,
			Defense:     defense,
			Speed:       speed,
			Special:     special,
			GifURL:      row[9],
			PNGURL:      row[10],
			Description: row[11],
		}

		pokemons = append(pokemons, pokemon)
	}

	return pokemons, nil
}
