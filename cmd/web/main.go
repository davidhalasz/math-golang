package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/davidhalasz/gomath/cmd/web/internal/config"
	"github.com/davidhalasz/gomath/cmd/web/internal/handlers"
	"github.com/davidhalasz/gomath/cmd/web/internal/helpers"
	"github.com/davidhalasz/gomath/cmd/web/internal/render"
)

const portnNumber = ":8080"

var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Starting application on port %s", portnNumber))

	srv := &http.Server{
		Addr:    portnNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")

	flag.Parse()

	app.InProduction = *inProduction
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create  template cache")
		return err
	}

	app.TemplateCache = tc

	handlers.NewHandlers(&app)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return nil
}
