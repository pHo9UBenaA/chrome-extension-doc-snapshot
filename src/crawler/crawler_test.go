package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/storage/mock"
)

func TestCrawler(t *testing.T) {
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
							<dt><a href="/api1">API 1</a></dt>
							<dt><a href="/api2">API 2</a></dt>
						</dl>
					</body>
				</html>
			`))
		case "/api1", "/api2":
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
	storage := mock.NewMockStorage()
	crawler := NewCrawler(storage)

	// Act
	// テストサーバーのURLをBaseURLとして使用
	doc, err := crawler.FetchAPIReference(ts.URL)
	if err != nil {
		t.Fatalf("FetchAPIReference failed: %v", err)
	}

	// Assert
	if doc == nil {
		t.Fatal("Expected document to be non-nil")
	}
}
