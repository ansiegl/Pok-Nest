package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/config"
	"github.com/ansiegl/Pok-Nest.git/internal/data"
	"github.com/ansiegl/Pok-Nest.git/internal/models"
	"github.com/ansiegl/Pok-Nest.git/internal/util"
	"github.com/ansiegl/Pok-Nest.git/internal/util/command"
	dbutil "github.com/ansiegl/Pok-Nest.git/internal/util/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var pokemonCSVPath string

func newSeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Inserts or updates fixtures to the database.",
		Long:  `Uses upsert to add test data to the current environment.`,
		Run: func(_ *cobra.Command, _ []string) {
			seedCmdFunc()
		},
	}

	cmd.Flags().StringVarP(&pokemonCSVPath, "file", "f", "/app/docs/pokemon.csv", "Path to CSV")

	return cmd
}

func seedCmdFunc() {
	err := command.WithServer(context.Background(), config.DefaultServiceConfigFromEnv(), func(ctx context.Context, s *api.Server) error {
		log := util.LogFromContext(ctx)

		err := ApplySeedFixtures(ctx, s.Config)
		if err != nil {
			log.Err(err).Msg("Error while applying seed fixtures")
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to apply migrations")
	}
}

func ApplySeedFixtures(ctx context.Context, config config.Server) error {
	log := util.LogFromContext(ctx)

	// Debugging: Ausgabe des CSV-Pfads
	log.Info().Str("csvPath", pokemonCSVPath).Msg("CSV Path received")

	// Überprüfe, ob die Datei existiert
	if _, err := os.Stat(pokemonCSVPath); os.IsNotExist(err) {
		return fmt.Errorf("CSV file does not exist: %s", pokemonCSVPath)
	}

	db, err := sql.Open("postgres", config.Database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	return dbutil.WithTransaction(ctx, db, func(tx boil.ContextExecutor) error {
		// Ausführlichere Debugging-Ausgabe
		pokemonCount, err := models.Pokemons().Count(ctx, tx)
		log.Info().
			Int64("pokemonCount", pokemonCount).
			Msg("Current Pokemon count in database")

		if err != nil {
			return fmt.Errorf("error checking existing Pokemon: %w", err)
		}

		if pokemonCount > 0 {
			log.Info().
				Int64("existingPokemonCount", pokemonCount).
				Msg("Pokemon data already exists in the database. Skipping seed.")
			return nil
		}

		fixtures := data.Upserts()

		for _, fixture := range fixtures {
			if err := fixture.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
				log.Error().Err(err).Msg("Failed to upsert fixture")
				return err
			}
		}

		// get pokemon data
		pokemons, err := data.LoadPokemonFromCSV(pokemonCSVPath)
		if err != nil {
			return fmt.Errorf("failed to load pokemon from CSV: %w", err)
		}

		for _, pokemon := range pokemons {
			if err := pokemon.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Error().Err(err).Msg("Failed to insert pokemon")
				return err
			}
		}

		log.Info().
			Int("fixturesCount", len(fixtures)).
			Int("pokemonCount", len(pokemons)).
			Msg("Successfully upserted fixtures and Pokemon")

		return nil
	})
}
