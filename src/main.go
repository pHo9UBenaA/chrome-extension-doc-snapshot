package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

const (
	SnapshotDir           = "__snapshot__"
	SnapshotFileExtension = ".md"
	BaseURL               = "https://developer.chrome.com/"
)

func ensureSnapshotDir() error {
	return os.MkdirAll(SnapshotDir, 0755)
}

func convertNodeToMarkdown(n *html.Node) (string, error) {
	htmlContent := renderNode(n)
	return htmltomarkdown.ConvertString(htmlContent)
}

func takeSnapshot(href string, markdown string) error {
	if err := ensureSnapshotDir(); err != nil {
		return err
	}

	filePath := path.Join(SnapshotDir, path.Base(href)+SnapshotFileExtension)

	return os.WriteFile(filePath, []byte(markdown), 0644)
}

func hasClass(n *html.Node, className string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			return strings.Contains(attr.Val, className)
		}
	}
	return false
}

func removeUnwantedElements(n *html.Node) {
	var next *html.Node

	for c := n.FirstChild; c != nil; c = next {
		next = c.NextSibling

		if c.Type != html.ElementNode {
			continue
		}

		switch {
		// パンくずリスト
		case c.Data == "div" && hasClass(c, "devsite-article-meta"):
				c.Parent.RemoveChild(c)
		// ヘッダーのツールチップ
		case c.Data == "h1" && hasClass(c, "devsite-page-title"):
			for inner := c.FirstChild; inner != nil; inner = inner.NextSibling {
				if inner.Data == "div" {
					inner.Parent.RemoveChild(inner)
				}
			}
		}
	}
}

func findArticle(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "article" {
		removeUnwantedElements(n)
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findArticle(c); result != nil {
			return result
		}
	}
	return nil
}

func renderNode(n *html.Node) string {
	buf := new(bytes.Buffer)
	html.Render(buf, n)
	return buf.String()
}

func extractArticle(doc *html.Node) (*html.Node, error) {
	if article := findArticle(doc); article != nil {
		return article, nil
	}
	return nil, fmt.Errorf("article not found")
}

func fetchAPIDetail(href string) (*html.Node, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Get(BaseURL + href)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return html.Parse(res.Body)
}

func findHrefInAnchor(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			for _, attr := range c.Attr {
				if attr.Key == "href" {
					return attr.Val
				}
			}
		}
	}
	return ""
}

func extractAPILinks(doc *html.Node) []string {
	var apiLinks []string

    var findLinks func(*html.Node)
    findLinks = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "dt" {
            if href := findHrefInAnchor(n); href != "" {
                apiLinks = append(apiLinks, href)
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            findLinks(c)
        }
    }

    findLinks(doc)
    return apiLinks
}

func fetchAPIReference() (*html.Node, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(BaseURL + "docs/extensions/reference/api")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return html.Parse(resp.Body)
}

func fetchAndSnapshotArticle(href string) error {
	doc, err := fetchAPIDetail(href)
	if err != nil {
		return err
	}

	articleNode, err := extractArticle(doc)
	if err != nil {
		return err
	}

	markdown, err := convertNodeToMarkdown(articleNode)
	if err != nil {
		return err
	}

	if err := takeSnapshot(href, markdown); err != nil {
		return err
	}

	return nil
}

func main() {
	log.Println("Start: Scrape and Snapshot")

	log.Println("Start: fetchAndSnapshotAPIReference")

	doc, err := fetchAPIReference()
	if err != nil {
		log.Fatalf("failed to fetchAPIReference: %v", err)
	}

	hrefList := extractAPILinks(doc)

	log.Println("Done:  fetchAndSnapshotAPIReference")

	var wg sync.WaitGroup
	var mu sync.Mutex
	errorChan := make(chan error, len(hrefList))

	for _, href := range hrefList {
		wg.Add(1)

		go func(href string) {
			defer wg.Done()

			baseName := path.Base(href)

			log.Printf("Start: fetchAndSnapshotArticle (%s)\n", baseName)

			if err := fetchAndSnapshotArticle(href); err != nil {
				mu.Lock()
				errorChan <- fmt.Errorf("failed to fetchAndSnapshotArticle (%s): %v", href, err)
				mu.Unlock()
				return
			}

			log.Printf("Done:  fetchAndSnapshotArticle (%s)\n", baseName)
		}(href)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	log.Println("Done: Scrape and Snapshot")
}
