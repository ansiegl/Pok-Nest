// Code generated by go-swagger; DO NOT EDIT.

package collection

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewDeletePokemonFromCollectionParams creates a new DeletePokemonFromCollectionParams object
// no default values defined in spec.
func NewDeletePokemonFromCollectionParams() DeletePokemonFromCollectionParams {

	return DeletePokemonFromCollectionParams{}
}

// DeletePokemonFromCollectionParams contains all the bound params for the delete pokemon from collection operation
// typically these are obtained from a http.Request
//
// swagger:parameters deletePokemonFromCollection
type DeletePokemonFromCollectionParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*pokemon's ID to delete
	  Required: true
	  In: path
	*/
	PokemonID string `param:"pokemonId"`
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeletePokemonFromCollectionParams() beforehand.
func (o *DeletePokemonFromCollectionParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rPokemonID, rhkPokemonID, _ := route.Params.GetOK("pokemonId")
	if err := o.bindPokemonID(rPokemonID, rhkPokemonID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DeletePokemonFromCollectionParams) Validate(formats strfmt.Registry) error {
	var res []error

	// pokemonId
	// Required: true
	// Parameter is provided by construction from the route

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindPokemonID binds and validates parameter PokemonID from path.
func (o *DeletePokemonFromCollectionParams) bindPokemonID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.PokemonID = raw

	return nil
}
