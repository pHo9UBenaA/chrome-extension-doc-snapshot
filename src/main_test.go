package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/storage"
)

func TestMain(t *testing.T) {
	// Arrange
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	os.Setenv("SNAPSHOT_DIR", tempDir)

	// テスト用のHTMLサーバーを立ち上げる
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/docs/extensions/reference/api":
			w.Write([]byte(`
				<!DOCTYPE html>
				<html>
					<body>
						<article>
							<dl>
								<dt><a href="/docs/extensions/reference/api/api1">API 1</a></dt>
								<dt><a href="/docs/extensions/reference/api/api2">API 2</a></dt>
							</dl>
						</article>
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

	// Act
	// テストサーバーのURLをBaseURLとして使用
	BaseURL = ts.URL

	// クローリングを実行
	main()

	// Assert
	// スナップショットが保存されたことを確認
	snapshotDir := filepath.Join(tempDir, storage.SnapshotDir)
	files, err := os.ReadDir(snapshotDir)
	if err != nil {
		t.Fatalf("Failed to read snapshot directory: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("Expected 2 snapshots, got %d", len(files))
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), storage.SnapshotFileExtension) {
			t.Errorf("Unexpected file extension: %s", file.Name())
		}
	}
}
