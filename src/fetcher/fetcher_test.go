package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/net/html"
)

func Test_FetchHTML(t *testing.T) {
	// Arrange
	// テスト用のHTMLサーバーを立ち上げる
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
				<head>
					<title>Test Page</title>
				</head>
				<body>
					<h1>Hello World</h1>
				</body>
			</html>
		`))
	}))
	defer ts.Close()

	// Act
	// テストサーバーのURLに対してリクエストを実行
	doc, err := FetchHTML(ts.URL, 30*time.Second)

	// Assert
	// エラーがないことを確認
	if err != nil {
		t.Fatalf("FetchHTML failed: %v", err)
	}

	// HTMLドキュメントが正しくパースされていることを確認
	if doc == nil {
		t.Fatal("Expected non-nil HTML document")
	}

	// タイトルが正しく取得できることを確認
	var title string
	var findTitle func(*html.Node)
	findTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTitle(c)
		}
	}
	findTitle(doc)

	if title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got '%s'", title)
	}
}

func Test_FetchHTML_Timeout(t *testing.T) {
	// Arrange
	// 意図的に遅延するテストサーバー
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("<html></html>"))
	}))
	defer ts.Close()

	// Act
	// 非常に短いタイムアウトを設定
	_, err := FetchHTML(ts.URL, 1*time.Millisecond)

	// Assert
	// タイムアウトエラーが発生することを確認
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}
}

func Test_FetchHTML_InvalidURL(t *testing.T) {
	// Arrange
	invalidURL := "http://invalid-url-that-does-not-exist"

	// Act
	_, err := FetchHTML(invalidURL, 30*time.Second)

	// Assert
	if err == nil {
		t.Fatal("Expected error for invalid URL, got nil")
	}
}
