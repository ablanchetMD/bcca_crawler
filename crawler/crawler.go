package crawler

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"fmt"
	"os"
	
	"golang.org/x/net/html"

)

func getTextContent(n *html.Node) string {
	var buf strings.Builder
	var extract func(*html.Node)
	extract = func(node *html.Node) {
		if node.Type == html.TextNode {
			buf.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(n)
	return strings.TrimSpace(buf.String())
}

func GetHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	if !strings.Contains(resp.Header.Get("content-type"), "text/html") {
		return "", errors.New("invalid content type : " + resp.Header.Get("content-type"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}


func GetURLsFromHTML(htmlBody string) ([]map[string]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	return_list := []map[string]string{}
	if err != nil {
		return nil, err
	}
	i := 0
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h4" {
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				if n.NextSibling != nil && n.NextSibling.Type == html.ElementNode && n.NextSibling.Data == "p" {
					pNode := n.NextSibling
					if pNode.NextSibling != nil && pNode.NextSibling.Type == html.ElementNode && pNode.NextSibling.Data == "ul" {
						ulNode := pNode.NextSibling
	
						// Extract text content
						h4Text := getTextContent(n)
						pText := getTextContent(pNode)
						ulText := getTextContent(ulNode)
	
						// Save to result
						return_list = append(return_list, map[string]string{
							"h4": h4Text,
							"p":  pText,
							"ul": ulText,
						})
						fmt.Fprintf(os.Stdout, "%d: %v : %v \n", i, h4Text,pText)

						i++
					}
				}
				
				
			}
			
			
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return return_list, nil
}