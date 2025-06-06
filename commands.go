package main

import (
	"bcca_crawler/ai_helper"
	"github.com/gorilla/mux"
	"bcca_crawler/api"
	"bcca_crawler/crawler"
	"bcca_crawler/internal/config"	
	"bcca_crawler/internal/auth"
	"bcca_crawler/routes"
	"time"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"io"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	library map[string]func(*config.Config, command) error
}



// func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.URL.Path == "/" {
// 			atomic.AddUint64(&cfg.fileserverHits, 1)
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func (c *commands) register(name string, f func(s *config.Config, c command) error) error {
	if c.library == nil {
		c.library = make(map[string]func(*config.Config, command) error) // Initialize the map if not already done
	}
	if _, exists := c.library[name]; exists {
		return errors.New("command already registered")
	}
	c.library[name] = f
	return nil
}

func (c *commands) run(s *config.Config, cmd command) error {
	f, ok := c.library[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func handlerGeoLocation(s *config.Config, cmd command) error {
	// Get the geolocation of a given IP address
	ip := cmd.Args[0]

	resp,err := api.GetIPGeoLocation(ip)
	// geo, err := s.GeoIPDB.City(net.ParseIP(ip))
	if err != nil {
		fmt.Println("Error getting geolocation: ", err)
		return err
	}
	fmt.Println("IP: ", resp)	
	return nil
}

func handlerAnalyzePDF(s *config.Config, cmd command) error {
	// Analyze a PDF
	err := ai_helper.GetAiData(s,cmd.Args[0])
	if err != nil {
		fmt.Println("Error analyzing PDF: ", err)
		return err 
	}
	
	return nil
}

func handlerSearchPubmed(s *config.Config, cmd command) error {

	url := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&term="+cmd.Args[0]+"&retmode=json"
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error:", err)
        return err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response:", err)
        return err
    }

    fmt.Println(string(body))
	return nil
}

func handlerCheckDatabase(s *config.Config, cmd command) error {
	// Check the database
	ctx := context.Background()
	payload,err := api.CMD_GetProtocolBy(s,ctx,"code",cmd.Args[0])
	if err != nil {
		fmt.Println("Error checking database: ", err)
		return err
	}
	fmt.Println("Payload: ")
	api.PrintStruct(payload)
	return nil
}

func handlerDeleteProtocol(s *config.Config, cmd command) error {
	// Check the database
	err := api.CMD_DeleteProtocol(s,cmd.Args[0])
	if err != nil {		
		return err
	}	
	return nil
}


func handlerCreateUser(s *config.Config, cmd command) error {
	// Create a new user
	email := cmd.Args[0]
	password := cmd.Args[1]

	user,err := api.HandleCLICreateUser(s, email, password)
	if err != nil {
		fmt.Println("Error creating user: ", err)
		return err
	}
	fmt.Println("User created successfully: ", user)
	return nil
}

func handlerCrawl(s *config.Config, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("missing URL argument")
	}
	// Crawl a website
	url := cmd.Args[0]
	
	html,err := crawler.GetHTML(url)
	if err != nil {
		fmt.Println("Error getting HTML: ", err)
	}
	_,err = crawler.GetURLsFromHTML(html)
	if err != nil {
		fmt.Println("Error getting URLs: ", err)
	}
	// fmt.Println("Links found: ", links)
	
	return nil
}


func handlerStartServer(s *config.Config, cmd command) error {
	// Start the server
	// Create a new instance of the server
	router := mux.NewRouter()
	routes.RegisterRoutes(router, s)
	// mux := http.NewServeMux()

	// routes.RegisterRoutes(mux, s)
	// Apply CORS middleware to the router
    corsRouter := auth.CORSMiddleware(router)
	
	// // Start the server//
	// wrappedMux := middleware.MiddlewareAuth(s,mux)
	// portString := "8080"
	srv := &http.Server{
		Addr:    ":" + s.ServerPort,
		Handler: corsRouter,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Server listening on port %s", s.ServerPort)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server: ", err)
		log.Fatal(err)
	}
	return nil
}
