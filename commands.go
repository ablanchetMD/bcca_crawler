package main

import (	
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
	library map[string]func(*ApiConfig, command) error
}



// func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.URL.Path == "/" {
// 			atomic.AddUint64(&cfg.fileserverHits, 1)
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func (c *commands) register(name string, f func(s *ApiConfig, c command) error) error {
	if c.library == nil {
		c.library = make(map[string]func(*ApiConfig, command) error) // Initialize the map if not already done
	}
	if _, exists := c.library[name]; exists {
		return errors.New("command already registered")
	}
	c.library[name] = f
	return nil
}

func (c *commands) run(s *ApiConfig, cmd command) error {
	f, ok := c.library[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}


func handlerStartServer(s *ApiConfig, cmd command) error {
	// Start the server
	// Create a new instance of the server
	mux := http.NewServeMux()

	
	// Define a handler
	mux.HandleFunc("GET /api/v1/protocols", func(w http.ResponseWriter, r *http.Request) {
		handleGetProtocols(s, w, r)		
	})	

	mux.HandleFunc("POST /api/v1/protocols", func(w http.ResponseWriter, r *http.Request) {
		handleCreateProtocol(s, w, r)		
	})

	// Start the server//
	//  wrappedMux := s.middlewareMetricsInc(mux)
	// portString := "8080"
	srv := &http.Server{
		Addr:    ":" + s.ServerPort,
		Handler: mux,
	}
	log.Printf("Server listening on port %s", s.ServerPort)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server: ", err)
		log.Fatal(err)
	}
	return nil
}
