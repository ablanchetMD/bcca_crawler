package config

import (
	"bcca_crawler/internal/database"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Db             *database.Queries
	Platform       string
	ServerPort	   string
	DatabaseUrl    string
	Secret         string
	Validate	   *validator.Validate
	
}