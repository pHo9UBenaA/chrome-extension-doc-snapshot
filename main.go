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
)

const (
	SnapshotDir           = "__snapshot__"
	SnapshotFileExtension = ".html"
	BaseURL               = "https://developer.chrome.com/"
)

func ensureSnapshotDir() error {
	if err := os.MkdirAll(SnapshotDir, 0755); err != nil {
		return fmt.Errorf("failed to make snapshot directory: %w", err)
	}

	return nil
}

func takeSnapshot(fileName string, content string) error {
	if err := ensureSnapshotDir(); err != nil {
		return fmt.Errorf("failed to ensure snapshot directory: %v", err)
	}

	filePath := path.Join(SnapshotDir, fileName+SnapshotFileExtension)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return nil
}

func findArticle(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "article" {
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
	var buf bytes.Buffer
	if err := html.Render(&buf, n); err != nil {
		return ""
	}
	return buf.String()
}

func scrapeArticle(href string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(BaseURL + href)
	if err != nil {
		return "", fmt.Errorf("failed to get article: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	article := findArticle(doc)
	if article == nil {
		return "", fmt.Errorf("article not found")
	}

	return renderNode(article), nil
}

func scrapeAPIReference() ([]string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(BaseURL + "docs/extensions/reference/api")
	if err != nil {
		return nil, fmt.Errorf("failed to get API reference: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var apiLinks []string
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "dt" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "a" {
					for _, attr := range c.Attr {
						if attr.Key == "href" {
							apiLinks = append(apiLinks, attr.Val)
							break
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}
	findLinks(doc)

	return apiLinks, nil
}

func scrapeAndSnapshotAPIReference() ([]string, error) {
	apiLinks, err := scrapeAPIReference()
	if err != nil {
		return nil, fmt.Errorf("failed to scrapeAPIReference: %w", err)
	}

	content := strings.Join(apiLinks, "\n")
	if err := takeSnapshot("api-reference", content); err != nil {
		return nil, fmt.Errorf("failed to takeSnapshot (api-reference): %w", err)
	}

	return apiLinks, nil
}

func scrapeAndSnapshotArticle(href string) error {
	content, err := scrapeArticle(href)
	if err != nil {
		return fmt.Errorf("failed to scrapeArticle: %w", err)
	}

	baseName := path.Base(href)
	if err := takeSnapshot(baseName, content); err != nil {
		return fmt.Errorf("failed to takeSnapshot (article): %w", err)
	}

	return nil
}

func main() {
	log.Println("Start: Scrape and Snapshot")

	log.Println("Start: scrapeAndSnapshotAPIReference")

	hrefList, err := scrapeAndSnapshotAPIReference()
	if err != nil {
		log.Fatalf("failed to scrapeAndSnapshotAPIReference: %v", err)
	}

	log.Println("Done:  scrapeAndSnapshotAPIReference")

	var wg sync.WaitGroup
	var mu sync.Mutex
	errorChan := make(chan error, len(hrefList))

	for _, href := range hrefList {
		wg.Add(1)

		go func(href string) {
			defer wg.Done()

			baseName := path.Base(href)

			log.Printf("Start: scrapeAndSnapshotArticle (%s)\n", baseName)

			if err := scrapeAndSnapshotArticle(href); err != nil {
				mu.Lock()
				errorChan <- fmt.Errorf("failed to scrapeAndSnapshotArticle (%s): %v", href, err)
				mu.Unlock()
				return
			}

			log.Printf("Done:  scrapeAndSnapshotArticle (%s)\n", baseName)
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

	// リクエスト間に少し待機を入れる
	time.Sleep(1 * time.Second)
}
