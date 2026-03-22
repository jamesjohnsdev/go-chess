// Package httpserver wires up the Huma-documented API for go-chess: health
// check now, live-game and chat endpoints (likely WebSocket) once
// internal/game grows session management. The frontend is a separate app
// that talks to this API.
package httpserver

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"github.com/jamesjohnsdev/go-chess/internal/game"
)

func New() http.Handler {
	router := chi.NewMux()
	api := huma.NewAPI(huma.DefaultConfig("go-chess", "0.1.0"), humachi.NewAdapter(router))
	store := game.NewStore()

	registerHealthz(api)
	registerCreateGame(api, store)
	registerGetGame(api, store)
	registerGameSocket(router, store)

	return router
}

type HealthzOutput struct {
	Body struct {
		Status string `json:"status" example:"ok"`
	}
}

func registerHealthz(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "healthz",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Summary:     "Health check",
	}, func(ctx context.Context, input *struct{}) (*HealthzOutput, error) {
		out := &HealthzOutput{}
		out.Body.Status = "ok"
		return out, nil
	})
}
