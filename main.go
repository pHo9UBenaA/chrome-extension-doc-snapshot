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
	if err := os.MkdirAll(SnapshotDir, 0755); err != nil {
		return fmt.Errorf("failed to make snapshot directory: %w", err)
	}

	return nil
}

func takeSnapshot(href string, content string) error {
	fileName := path.Base(href)

	if err := ensureSnapshotDir(); err != nil {
		return fmt.Errorf("failed to ensure snapshot directory: %v", err)
	}

	filePath := path.Join(SnapshotDir, fileName+SnapshotFileExtension)

	markdown, err := htmltomarkdown.ConvertString(content)
	if err != nil {
		return fmt.Errorf("failed to convert HTML to Markdown: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func hasClass(n *html.Node, className string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" && strings.Contains(attr.Val, className) {
			return true
		}
	}
	return false
}

func removeUnwantedElements(n *html.Node) {
	var next *html.Node

	for c := n.FirstChild; c != nil; c = next {
		next = c.NextSibling

		if c.Type == html.ElementNode {
			// パンくずリスト
			if c.Data == "div" && hasClass(c, "devsite-article-meta") {
				c.Parent.RemoveChild(c)
			}

			// ヘッダーのツールチップ
			if c.Data == "h1" && hasClass(c, "devsite-page-title") {
				for inner := c.FirstChild; inner != nil; inner = inner.NextSibling {
					if inner.Data == "div" {
						inner.Parent.RemoveChild(inner)
					}
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
	var buf bytes.Buffer
	if err := html.Render(&buf, n); err != nil {
		return ""
	}
	return buf.String()
}

func fetchArticle(href string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Get(BaseURL + href)
	if err != nil {
		return "", fmt.Errorf("failed to get article: %w", err)
	}
	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	article := findArticle(doc)
	if article == nil {
		return "", fmt.Errorf("article not found")
	}

	return renderNode(article), nil
}

func fetchAPIReference() ([]string, error) {
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

func fetchAndSnapshotAPIReference() ([]string, error) {
	apiLinks, err := fetchAPIReference()
	if err != nil {
		return nil, fmt.Errorf("failed to fetchAPIReference: %w", err)
	}

	content := strings.Join(apiLinks, "\n")
	
	if err := takeSnapshot("api-reference", content); err != nil {
		return nil, fmt.Errorf("failed to takeSnapshot (api-reference): %w", err)
	}

	return apiLinks, nil
}

func fetchAndSnapshotArticle(href string) error {
	content, err := fetchArticle(href)
	if err != nil {
		return fmt.Errorf("failed to fetchArticle: %w", err)
	}

	if err := takeSnapshot(href, content); err != nil {
		return fmt.Errorf("failed to takeSnapshot (article): %w", err)
	}

	return nil
}

func main() {
	log.Println("Start: Scrape and Snapshot")

	log.Println("Start: fetchAndSnapshotAPIReference")

	hrefList, err := fetchAndSnapshotAPIReference()
	if err != nil {
		log.Fatalf("failed to fetchAndSnapshotAPIReference: %v", err)
	}

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

	// リクエスト間に少し待機を入れる
	time.Sleep(1 * time.Second)
}
