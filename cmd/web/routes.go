package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/home", app.HomePage)
	mux.Get("/statistics", app.StatisticsPage)
	mux.Get("/statistics/mean", app.Mean)
	mux.Get("/statistics/median", app.Median)
	mux.Get("/statistics/std-deviation-variance", app.StdVar)
	mux.Get("/statistics/pdf", app.PDF)
	mux.Get("/statistics/binomial", app.Binomial)
	mux.Get("/statistics/poisson", app.Poisson)
	mux.Get("/statistics/covcor", app.CovCor)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
