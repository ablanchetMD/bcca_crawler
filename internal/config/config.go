package config

import (
	"bcca_crawler/internal/database"
	"github.com/go-playground/validator/v10"
	"database/sql"
)

type Config struct {
	Db             *database.Queries
	Database	   *sql.DB
	Platform       string
	ServerPort	   string
	DatabaseUrl    string
	Secret         string
	GeminiApiKey   string
	MailGunApiKey  string
	Validate	   *validator.Validate
	
}