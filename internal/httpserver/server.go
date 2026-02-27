// Package httpserver wires up the HTTP handlers for the web version of the
// game: static assets and a Huma-documented API now, live-game and chat
// endpoints (likely WebSocket) once internal/game grows session management.
package httpserver

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func New() http.Handler {
	router := chi.NewMux()
	api := huma.NewAPI(huma.DefaultConfig("go-chess", "0.1.0"), humachi.NewAdapter(router))

	registerHealthz(api)

	router.Handle("/*", http.FileServer(http.Dir("web/static")))
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
