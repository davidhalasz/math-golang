package main

import (
	"net/http"

	"github.com/davidhalasz/gomath/cmd/web/internal/config"
	"github.com/davidhalasz/gomath/cmd/web/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", handlers.HomePage)
	mux.Get("/statistics", handlers.StatisticsPage)
	mux.Get("/statistics/mean", handlers.Mean)
	mux.Get("/statistics/median", handlers.Median)
	mux.Get("/statistics/std-deviation-variance", handlers.StdVar)
	mux.Get("/statistics/pdf", handlers.PDF)
	mux.Get("/statistics/binomial", handlers.Binomial)
	mux.Get("/statistics/poisson", handlers.Poisson)
	mux.Get("/statistics/covcor", handlers.CovCor)
	mux.Get("/statistics/linear-regression", handlers.LinearRegression)

	mux.Get("/ai-basics", handlers.AiPage)
	mux.Get("/ai-basics/b", handlers.CallDLS)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
