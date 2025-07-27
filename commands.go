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
	"strings"
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
	if len(cmd.Args) < 1 {
		return errors.New("missing PDF URLs argument")
	}
	// Validate all URLs
    for _, url := range cmd.Args {
        if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
            return fmt.Errorf("invalid URL: %s. URL must start with http:// or https://", url)
        }
    }
	// Call the AI helper to analyze the PDF
	// This function will handle the PDF analysis and return any errors encountered

	ai_helper.RunAllLinks(s,cmd.Args)	
	
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

func handlerResetDatabase(s *config.Config, cmd command) error {
	// Reset the database
	err := api.CMD_ResetDatabase(s)
	if err != nil {
		fmt.Println("Error resetting database: ", err)
		return err
	}
	fmt.Println("Database reset successfully")
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
	WebProtocol,err := crawler.GetURLsFromHTML(html)
	if err != nil {
		fmt.Println("Error getting URLs: ", err)
	}

	protocol_list := make([]string, 0)

	for _, protocol := range WebProtocol {
		fmt.Println("Protocol Name: ", protocol.Code)		
		for _, link:= range protocol.Links {
			textVal, ok := link["text"]
			if ok && strings.Contains(textVal, "Protocol") {
				href := link["href"]
				fmt.Println("➡ Found Protocol link:", href)
				protocol_list = append(protocol_list, href)

			}
		}		
	}	
	
	ai_helper.RunAllLinks(s,protocol_list)
	
	return nil
}

func handlerSingleCrawl(s *config.Config, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("missing URL argument")
	}
	// Get a single protocol from a website
			
	ai_helper.RunAllLinks(s,cmd.Args)
	
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
