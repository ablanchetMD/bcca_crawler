package main

import (
	"github.com/joho/godotenv"
	"os"
	"fmt"
	"bcca_crawler/internal/database"
	_ "github.com/lib/pq"
	
)

type ApiConfig struct {
	Db             *database.Queries
	Platform       string
	ServerPort		   string
	DatabaseUrl    string
	Secret         string
	fileserverHits uint64
	
}

func main() {
	cfg := &ApiConfig{}
	godotenv.Load(".env")
	cfg.Platform = os.Getenv("PLATFORM")
	cfg.Secret = os.Getenv("SECRET")
	cfg.DatabaseUrl = os.Getenv("DB_URL")
	cfg.ServerPort = os.Getenv("PORT")
	commands := commands{}
	commands.register("serve", handlerStartServer)
	commands.register("readPdf", handlerReadPdf)

	//http://www.bccancer.bc.ca/health-professionals/clinical-resources/chemotherapy-protocols/lymphoma-myeloma

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No command provided")
		os.Exit(1)
		
	}

	err := commands.run(cfg, command{Name: args[0], Args: args[1:]})
	if err != nil {
		fmt.Println("Error running command :", err)
		os.Exit(1)
	}
	
}