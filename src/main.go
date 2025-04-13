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

var (
	baseURL = "https://developer.chrome.com/"
)

const (
	apiReferencePath       = "/docs/extensions/reference/api"
	maxConcurrentDownloads = 10
)

type DocumentProcessError struct {
	documentPath string
	err          error
}

func (e *DocumentProcessError) Error() string {
	return fmt.Sprintf("Document: '%s'\nError: %v", e.documentPath, e.err)
}

func getFilenameFromPath(documentPath string) (string, error) {
	// ex. ["", "/tabs"], ["", "/devtools/performance"]
	parts := strings.Split(documentPath, apiReferencePath)

	if parts[0] != "" || len(parts) != 2 {
		return "", fmt.Errorf("Failed to get filename from path: %s", parts)
	}

	path := strings.TrimPrefix(parts[1], "/")
	return strings.ReplaceAll(path, "/", "-"), nil
}

func fetchAPIDocumentLinks() ([]string, error) {
	referencePageURL := baseURL + apiReferencePath
	doc, err := crawler.FetchHTML(referencePageURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch API reference page: %w", err)
	}

	links, err := parser.ExtractAPILinks(doc)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract API links: %w", err)
	}

	return links, nil
}

func snapshotDocument(documentPath string) error {
	// ドキュメントのダウンロード
	fullURL := baseURL + documentPath
	doc, err := crawler.FetchHTML(fullURL)
	if err != nil {
		return &DocumentProcessError{
			documentPath: documentPath,
			err:          fmt.Errorf("Failed to fetch HTML: %w", err),
		}
	}

	// 記事本文の抽出
	article, err := parser.ExtractArticle(doc)
	if err != nil {
		return &DocumentProcessError{
			documentPath: documentPath,
			err:          fmt.Errorf("Failed to extract article: %w", err),
		}
	}

	// Markdownへの変換
	markdown, err := converter.NodeToMarkdown(article)
	if err != nil {
		return &DocumentProcessError{
			documentPath: documentPath,
			err:          fmt.Errorf("Failed to convert to Markdown: %w", err),
		}
	}

	// スナップショットの保存
	filename, err := getFilenameFromPath(documentPath)
	if err != nil {
		return &DocumentProcessError{
			documentPath: documentPath,
			err:          fmt.Errorf("Failed to get filename: %w", err),
		}
	}

	if err := storage.TakeSnapshot(filename, markdown); err != nil {
		return &DocumentProcessError{
			documentPath: documentPath,
			err:          fmt.Errorf("Failed to save snapshot: %w", err),
		}
	}

	log.Printf("✅ Document '%s' processed successfully", documentPath)
	return nil
}

func processDocumentsConcurrently(documentPaths []string) error {
	workerLimit := make(chan struct{}, maxConcurrentDownloads)
	errors := make(chan error, len(documentPaths))
	var wg sync.WaitGroup

	// ドキュメント処理の実行
	for _, path := range documentPaths {
		wg.Add(1)
		go func(docPath string) {
			defer wg.Done()
			workerLimit <- struct{}{} // 同時実行数の制限
			defer func() { <-workerLimit }()

			if err := snapshotDocument(docPath); err != nil {
				errors <- err
			}
		}(path)
	}

	// エラー収集
	var processErrors []error
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// 全ての処理が完了するまでエラーを収集
	go func() {
		<-done
		close(errors)
	}()

	for err := range errors {
		processErrors = append(processErrors, err)
	}

	if len(processErrors) > 0 {
		return fmt.Errorf("%d documents processing failed: %v", len(processErrors), processErrors)
	}

	return nil
}

func main() {
	log.Println("🚀 Starting snapshot processing for Chrome Extension API documents")

	documentPaths, err := fetchAPIDocumentLinks()
	if err != nil {
		log.Fatalf("❌ Failed to fetch API document links: %v", err)
	}
	log.Printf("📝 %d API documents detected", len(documentPaths))

	if err := processDocumentsConcurrently(documentPaths); err != nil {
		log.Fatalf("❌ Failed to process documents: %v", err)
	}

	log.Println("✨ All documents processed successfully")
}
