package config

import (
	"html/template"
	"log"
)

type Config struct {
	Port int
	Env  string
	Api  string
}

type AppConfig struct {
	UseCache      bool
	Config        Config
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	TemplateCache map[string]*template.Template
	Version       string
	InProduction  bool
}
