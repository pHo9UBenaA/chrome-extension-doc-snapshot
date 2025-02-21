package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	SnapshotDir = "__snapshot__"
	SnapshotFileExtension = ".txt"
	BaseURL     = "https://developer.chrome.com/"
)

func prepareCollector() *colly.Collector {
	return colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(2),
	)
}

func ensureSnapshotDir() error {
	if err := os.MkdirAll(SnapshotDir, 0755); err != nil {
		return fmt.Errorf("failed to make snapshot directory: %w", err)
	}

	return nil
}

func takeSnapshot(fileName string, content string) error {
	if err := ensureSnapshotDir(); err != nil {
		log.Fatalf("failed to ensure snapshot directory: %v", err)
	}

	filePath := path.Join(SnapshotDir, fileName+SnapshotFileExtension)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return nil
}

func scrapeAPIReference() ([]string, error) {
	c := prepareCollector()

	var apiLinks []string

	c.OnHTML("dl", func(e *colly.HTMLElement) {
		e.ForEach("dt", func(_ int, el *colly.HTMLElement) {
			href := el.ChildAttr("a", "href")
			if href != "" {
				apiLinks = append(apiLinks, href)
			}
		})
	})

	if err := c.Visit(BaseURL + "docs/extensions/reference/api"); err != nil {
		return nil, fmt.Errorf("failed to visit API reference: %w", err)
	}

	return apiLinks, nil
}

func scrapeArticle(href string) (string, error) {
	c := prepareCollector()
	var articleContent string
	c.OnHTML("article", func(e *colly.HTMLElement) {
		articleContent = e.Text
	})

	if err := c.Visit(BaseURL + href); err != nil {
		return "", fmt.Errorf("failed to visit article: %w", err)
	}

	return articleContent, nil
}

func scrapeAndSnapshotAPIReference() ([]string, error) {
	apiLinks, err := scrapeAPIReference()
	if err != nil {
		return nil, fmt.Errorf("failed to scrapeAPIReference: %w", err)
	}

	content := strings.Join(apiLinks, "\n")
	if err := takeSnapshot("index", content); err != nil {
		return nil, fmt.Errorf("failed to takeSnapshot (index): %w", err)
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

	for _, href := range hrefList {
		baseName := path.Base(href)

		log.Printf("Start: scrapeAndSnapshotArticle (%s)\n", baseName)

		if err := scrapeAndSnapshotArticle(href); err != nil {
			log.Fatalf("failed to scrapeAndSnapshotArticle (%s): %v", href, err)
			continue
		}

		log.Printf("Done:  scrapeAndSnapshotArticle (%s)\n", baseName)

		time.Sleep(1 * time.Second)
	}

	log.Println("Done: Scrape and Snapshot")
}
