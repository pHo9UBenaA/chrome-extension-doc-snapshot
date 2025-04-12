package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/converter"
	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/crawler"
	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/parser"
	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/storage"
)

// BaseURLはクロール対象のベースURL
var BaseURL = "https://developer.chrome.com/"

// APIReferencePathはChrome Extension APIリファレンスのパス
const APIReferencePath = "/docs/extensions/reference/api"

func extractAPILinks() ([]string, error) {
	doc, err := crawler.FetchHTML(BaseURL + APIReferencePath)
	if err != nil {
		return nil, err
	}
	return parser.ExtractAPILinks(doc)
}

func snapshotArticle(href string) error {
	doc, err := crawler.FetchHTML(BaseURL + href)
	if err != nil {
		return err
	}

	// 記事を抽出
	articleNode, err := parser.ExtractArticle(doc)
	if err != nil {
		return err
	}

	// Markdownに変換
	markdown, err := converter.ConvertNodeToMarkdown(articleNode)
	if err != nil {
		return err
	}

	// URLからパス部分を抽出し、最後の部分だけを取得
	path := href
	if idx := strings.LastIndex(path, "/"); idx != -1 {
		path = path[idx+1:]
	}

	// スナップショットとして保存
	if err := storage.TakeSnapshot(path, markdown); err != nil {
		return fmt.Errorf("failed to take snapshot: %v", err)
	}

	return nil
}

func main() {
	log.Println("Start: Scrape and Snapshot")

	// ストレージの初期化
	storage.EnsureSnapshotDir()

	log.Println("Scrape: API Reference")

	hrefList, err := extractAPILinks()
	if err != nil {
		log.Fatalf("failed to extractAPILinks: %v", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errorChan := make(chan error, len(hrefList))

	for _, href := range hrefList {
		wg.Add(1)

		go func(href string) {
			defer wg.Done()

			log.Printf("Scrape: Article (%s)", href)

			if err := snapshotArticle(href); err != nil {
				mu.Lock()
				errorChan <- fmt.Errorf("failed to snapshotArticle (%s): %v", href, err)
				mu.Unlock()
				return
			}
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
