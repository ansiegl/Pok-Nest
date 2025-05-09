// Code generated by go run -tags scripts scripts/handlers/gen_handlers.go; DO NOT EDIT.
package handlers

import (
	"github.com/ansiegl/Pok-Nest.git/internal/api"
	"github.com/ansiegl/Pok-Nest.git/internal/api/handlers/auth"
	"github.com/ansiegl/Pok-Nest.git/internal/api/handlers/collection"
	"github.com/ansiegl/Pok-Nest.git/internal/api/handlers/common"
	"github.com/ansiegl/Pok-Nest.git/internal/api/handlers/pokemon"
	"github.com/ansiegl/Pok-Nest.git/internal/api/handlers/push"
	"github.com/labstack/echo/v4"
)

func AttachAllRoutes(s *api.Server) {
	// attach our routes
	s.Router.Routes = []*echo.Route{
		auth.GetUserInfoRoute(s),
		auth.PostChangePasswordRoute(s),
		auth.PostForgotPasswordCompleteRoute(s),
		auth.PostForgotPasswordRoute(s),
		auth.PostLoginRoute(s),
		auth.PostLogoutRoute(s),
		auth.PostRefreshRoute(s),
		auth.PostRegisterRoute(s),
		collection.DeletePokemonRoute(s),
		collection.GetCollectionPokemonDetailRoute(s),
		collection.GetCollectionPokemonRoute(s),
		collection.PostAddPokemonToCollectionRoute(s),
		collection.PostSearchPokemonInCollectionRoute(s),
		collection.PutEditPokemonRoute(s),
		common.GetHealthyRoute(s),
		common.GetReadyRoute(s),
		common.GetSwaggerRoute(s),
		common.GetVersionRoute(s),
		pokemon.GetPokemonDetailsRoute(s),
		pokemon.GetPokemonDownloadRoute(s),
		pokemon.GetPokemonRoute(s),
		pokemon.PostSearchPokemonRoute(s),
		push.GetPushTestRoute(s),
		push.PostUpdatePushTokenRoute(s),
	}
}
