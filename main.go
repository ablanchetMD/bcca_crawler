package main

import (
	"github.com/joho/godotenv"
	"os"
	"fmt"
	"bcca_crawler/internal/database"
	_ "github.com/lib/pq"
	"database/sql"
	"github.com/go-playground/validator/v10"
	
)

type ApiConfig struct {
	Db             *database.Queries
	Platform       string
	ServerPort		   string
	DatabaseUrl    string
	Secret         string	
	
}
var validate *validator.Validate

func init() {
	godotenv.Load(".env")
	validate = validator.New()
	validate.RegisterValidation("tumorgroup", tumorGroupValidator)
}

func main() {
	cfg := &ApiConfig{}	
	cfg.Platform = os.Getenv("PLATFORM")
	cfg.Secret = os.Getenv("SECRET")
	cfg.DatabaseUrl = os.Getenv("DB_URL")
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		fmt.Println("Error fetching database: ", err)
		return
	}
	defer db.Close()
	dbQueries := database.New(db)
	cfg.Db = dbQueries
	cfg.ServerPort = os.Getenv("PORT")
	commands := commands{}
	commands.register("serve", handlerStartServer)
	

	//http://www.bccancer.bc.ca/health-professionals/clinical-resources/chemotherapy-protocols/lymphoma-myeloma

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	err = commands.run(cfg, command{Name: args[0], Args: args[1:]})
	if err != nil {
		fmt.Println("Error running command :", err)
		os.Exit(1)
	}
	
}