package main

import (
	"encoding/xml"
	"strings"
	"fmt"
	"io"
	"net/http"
	"os"
	"bytes"	
	
)

type Document struct {
	Title  string   `xml:"head>title"`
	Meta   []Meta   `xml:"head>meta"`
	Pages  []Page   `xml:"body>div"`
}

type Meta struct {
	Name    string `xml:"name,attr"`
	Content string `xml:"content,attr"`
}

type Page struct {
	Paragraphs []string   `xml:"p"`
	Links   []string `xml:"div.annotation>a"`
}

func ParseXMLDocument(xmlData string) (*Document, error) {

	  xmlData = strings.Replace(xmlData, `<?xml version="1.1" encoding="UTF-8"?>`, "", 1)
	var doc Document
	err := xml.Unmarshal([]byte(xmlData), &doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// ExtractText sends a PDF to Tika Server and retrieves the text content
func ExtractText(filePath string, tikaURL string) (string, error) {
	// Open the PDF file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	req,err := http.NewRequest("PUT", tikaURL+"/tika", bytes.NewReader(fileBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/pdf")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()	

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Read the response body
	responseText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(responseText), nil
}

// Helper function to trim newlines and extra spaces in the content
func CleanText(text string) string {
	// Trim leading and trailing whitespaces (including newlines)
	// text = strings.TrimSpace(text)

	// Replace any newlines within the text with a space or remove them
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")

	// Optionally, you can collapse multiple spaces to a single one
	text = strings.Join(strings.Fields(text), " ")

	return text
}