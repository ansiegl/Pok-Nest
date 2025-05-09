// Code generated by SQLBoiler 4.18.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"strconv"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/strmangle"
)

// M type is for providing columns and column values to UpdateAll.
type M map[string]interface{}

// ErrSyncFail occurs during insert when the record could not be retrieved in
// order to populate default value information. This usually happens when LastInsertId
// fails or there was a primary key configuration that was not resolvable.
var ErrSyncFail = errors.New("models: failed to synchronize data after insert")

type insertCache struct {
	query        string
	retQuery     string
	valueMapping []uint64
	retMapping   []uint64
}

type updateCache struct {
	query        string
	valueMapping []uint64
}

func makeCacheKey(cols boil.Columns, nzDefaults []string) string {
	buf := strmangle.GetBuffer()

	buf.WriteString(strconv.Itoa(cols.Kind))
	for _, w := range cols.Cols {
		buf.WriteString(w)
	}

	if len(nzDefaults) != 0 {
		buf.WriteByte('.')
	}
	for _, nz := range nzDefaults {
		buf.WriteString(nz)
	}

	str := buf.String()
	strmangle.PutBuffer(buf)
	return str
}

// Enum values for PokemonType
const (
	PokemonTypeFire     string = "Fire"
	PokemonTypeWater    string = "Water"
	PokemonTypeGrass    string = "Grass"
	PokemonTypeElectric string = "Electric"
	PokemonTypePsychic  string = "Psychic"
	PokemonTypeGhost    string = "Ghost"
	PokemonTypeDragon   string = "Dragon"
	PokemonTypeFairy    string = "Fairy"
	PokemonTypeFlying   string = "Flying"
	PokemonTypeDark     string = "Dark"
	PokemonTypeFighting string = "Fighting"
	PokemonTypeBug      string = "Bug"
	PokemonTypeNormal   string = "Normal"
	PokemonTypeRock     string = "Rock"
	PokemonTypeGround   string = "Ground"
	PokemonTypePoison   string = "Poison"
	PokemonTypeSteel    string = "Steel"
	PokemonTypeIce      string = "Ice"
)

func AllPokemonType() []string {
	return []string{
		PokemonTypeFire,
		PokemonTypeWater,
		PokemonTypeGrass,
		PokemonTypeElectric,
		PokemonTypePsychic,
		PokemonTypeGhost,
		PokemonTypeDragon,
		PokemonTypeFairy,
		PokemonTypeFlying,
		PokemonTypeDark,
		PokemonTypeFighting,
		PokemonTypeBug,
		PokemonTypeNormal,
		PokemonTypeRock,
		PokemonTypeGround,
		PokemonTypePoison,
		PokemonTypeSteel,
		PokemonTypeIce,
	}
}

// Enum values for ProviderType
const (
	ProviderTypeFCM string = "fcm"
	ProviderTypeApn string = "apn"
)

func AllProviderType() []string {
	return []string{
		ProviderTypeFCM,
		ProviderTypeApn,
	}
}
