package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	SnapshotDir           = "__snapshot__"
	SnapshotFileExtension = ".txt"
	BaseURL               = "https://developer.chrome.com/"
)

var globalCollector *colly.Collector

func init() {
	globalCollector = colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(1),
	)

	globalCollector.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Delay:      1 * time.Second,
		// 適当な数。特になくても弾かれなかった。
		Parallelism: 5,
	})
}

func prepareCollector() *colly.Collector {
	return globalCollector.Clone()
}

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

func scrapeAPIReference() ([]string, error) {
	c := prepareCollector()

	var apiLinks []string
	var mu sync.Mutex

	c.OnHTML("dl", func(e *colly.HTMLElement) {
		e.ForEach("dt", func(_ int, el *colly.HTMLElement) {
			href := el.ChildAttr("a", "href")
			if href != "" {
				mu.Lock()
				apiLinks = append(apiLinks, href)
				mu.Unlock()
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
	var mu sync.Mutex

	c.OnHTML("article", func(e *colly.HTMLElement) {
		mu.Lock()
		articleContent = e.Text
		mu.Unlock()
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
}
