package main

import (
	"bcca_crawler/api"
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/middleware"
	"bcca_crawler/routes"
	"errors"
	"fmt"
	"log"
	"net/http"
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


func handlerStartServer(s *config.Config, cmd command) error {
	// Start the server
	// Create a new instance of the server
	mux := http.NewServeMux()

	routes.RegisterRoutes(mux, s)
	
	// Start the server//
	wrappedMux := middleware.MiddlewareAuth(s,mux)
	// portString := "8080"
	srv := &http.Server{
		Addr:    ":" + s.ServerPort,
		Handler: wrappedMux,
	}
	log.Printf("Server listening on port %s", s.ServerPort)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server: ", err)
		log.Fatal(err)
	}
	return nil
}
