package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/crawler"
	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/parser"
	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/storage"
)

var (
	// BaseURLはクロール対象のベースURL
	BaseURL = "https://developer.chrome.com/"
)

func main() {
	log.Println("Start: Scrape and Snapshot")

	// ストレージの初期化
	storage := storage.NewFileStorage()
	crawler := crawler.NewCrawler(storage)

	log.Println("Start: fetchAndSnapshotAPIReference")

	doc, err := crawler.FetchAPIReference(BaseURL)
	if err != nil {
		log.Fatalf("failed to fetchAPIReference: %v", err)
	}

	hrefList := parser.ExtractAPILinks(doc)

	log.Println("Done:  fetchAndSnapshotAPIReference")

	var wg sync.WaitGroup
	var mu sync.Mutex
	errorChan := make(chan error, len(hrefList))

	for _, href := range hrefList {
		wg.Add(1)

		go func(href string) {
			defer wg.Done()

			log.Printf("Start: fetchAndSnapshotArticle (%s)\n", href)

			if err := crawler.FetchAndSnapshotArticle(BaseURL + href); err != nil {
				mu.Lock()
				errorChan <- fmt.Errorf("failed to fetchAndSnapshotArticle (%s): %v", href, err)
				mu.Unlock()
				return
			}

			log.Printf("Done:  fetchAndSnapshotArticle (%s)\n", href)
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
