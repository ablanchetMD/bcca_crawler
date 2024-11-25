package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	library map[string]func(*ApiConfig, command) error
}

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w http.ResponseWriter, name string, data map[string]interface{}, c context.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	tmpl := template.Must(template.ParseGlob("views/*.html"))
	return &Templates{tmpl}
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			atomic.AddUint64(&cfg.fileserverHits, 1)
		}
		next.ServeHTTP(w, r)
	})
}

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

func handlerReadPdf(s *ApiConfig, cmd command) error {
	wdirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory: ", err)
		log.Fatal(err)
	}
	pdfFile := "/assets/pdf/LYACAL_Protocol.pdf"
	pdfPath := wdirectory + pdfFile
	fmt.Println("pdfPath: ", pdfPath)
	content, err := ExtractText(pdfPath, "http://localhost:9998")
	if err != nil {
		fmt.Println("Error reading pdf: ", err)
		log.Fatal(err)
	}

	doc, err := ParseXMLDocument(content)
	if err != nil {
		log.Fatalf("Failed to parse XML: %v", err)
	}

	// Print metadata
	fmt.Println("Title:", doc.Title)
	for _, meta := range doc.Meta {
		fmt.Printf("Meta: Name=%s, Content=%s\n", meta.Name, meta.Content)
	}

	// Print page content and links
	for i, page := range doc.Pages {
		for _, paragraph := range page.Paragraphs {

			fmt.Printf("%s\n", CleanText(paragraph))
		}

		for _, link := range page.Links {
			fmt.Printf("Page %d Link: %s\n", i+1, link)
		}
	}

	return nil
}

func handlerStartServer(s *ApiConfig, cmd command) error {
	// Start the server
	// Create a new instance of the server
	mux := http.NewServeMux()

	tmpl := newTemplate()

	// Define a handler
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		raw,err := s.Db.GetProtocols(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error fetching protocols")
			return
		}
		var response []Protocol
		for _, protocol := range raw {
			response = append(response, mapProtocolStruct(protocol))
		}

		data := map[string]interface{}{
			"protocols": response,
		}

		// Execute the template with data
		if err := tmpl.Render(w, "index", data, r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Define a handler

	mux.HandleFunc("POST /count", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&s.fileserverHits, 1)
		data := map[string]interface{}{
			"Title": "Home",
			"Count": fmt.Sprintf("%d", s.fileserverHits),
		}

		// Execute the template with data
		if err := tmpl.Render(w, "count", data, r.Context()); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
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
