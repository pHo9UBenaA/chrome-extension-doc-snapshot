package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/crawler"
	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/parser"
	storage_mock "github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/storage/mock"
)

func TestMain(t *testing.T) {
	// Arrange
	// テスト用のHTMLサーバーを立ち上げる
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/docs/extensions/reference/api":
			w.Write([]byte(`
				<!DOCTYPE html>
				<html>
					<body>
						<dl>
							<dt><a href="/docs/extensions/reference/api/api1">API 1</a></dt>
							<dt><a href="/docs/extensions/reference/api/api2">API 2</a></dt>
						</dl>
					</body>
				</html>
			`))
		case "/docs/extensions/reference/api/api1", "/docs/extensions/reference/api/api2":
			w.Write([]byte(`
				<!DOCTYPE html>
				<html>
					<body>
						<article>
							<h1>Test Article</h1>
							<p>This is a test article.</p>
						</article>
					</body>
				</html>
			`))
		}
	}))
	defer ts.Close()

	// モックストレージを使用
	storage := storage_mock.NewMockStorage()
	crawler := crawler.NewCrawler(storage)

	// Act
	// テストサーバーのURLをBaseURLとして使用
	BaseURL := ts.URL

	// クローリングを実行
	doc, err := crawler.FetchAPIReference(BaseURL)
	if err != nil {
		t.Fatalf("FetchAPIReference failed: %v", err)
	}

	hrefList, err := parser.ExtractAPILinks(doc)
	if err != nil {
		t.Fatalf("ExtractAPILinks failed: %v", err)
	}

	for _, href := range hrefList {
		if err := crawler.FetchAndSnapshotArticle(ts.URL + href); err != nil {
			t.Fatalf("FetchAndSnapshotArticle failed: %v", err)
		}
	}

	// Assert
	// スナップショットが保存されたことを確認
	if len(storage.Snapshots) != 2 {
		t.Errorf("Expected 2 snapshots, got %d", len(storage.Snapshots))
	}

	// 各スナップショットの内容を確認
	for _, href := range hrefList {
		// ファイル名を生成（最後の部分を使用）
		fileName := href
		if idx := strings.LastIndex(href, "/"); idx != -1 {
			fileName = href[idx+1:]
		}

		content, exists := storage.Snapshots[fileName]
		if !exists {
			t.Errorf("Snapshot for %s was not saved", href)
			continue
		}

		if !strings.Contains(content, "Test Article") {
			t.Errorf("Snapshot for %s does not contain expected content", href)
		}
	}
}
