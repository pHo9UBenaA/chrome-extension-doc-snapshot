package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
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
								<dt><a href="/docs/extensions/reference/api/api">API 1</a></dt>
								<dt><a href="/docs/extensions/reference/api/enterprise/api">API 2</a></dt>
							</dl>
						</article>
					</body>
				</html>
			`))
		case "/docs/extensions/reference/api/api", "/docs/extensions/reference/api/enterprise/api":
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
	baseURL = ts.URL

	// クローリングを実行
	main()

	// Assert
	// スナップショットが保存されたことを確認
	snapshotDir := storage.GetSnapshotDirPath()
	files, err := os.ReadDir(snapshotDir)
	if err != nil {
		t.Fatalf("Failed to read snapshot directory: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("Expected 2 snapshots, got %d", len(files))
	}

	// スナップショットのファイル名を確認
	expectedFilenames := []string{"api.md", "enterprise-api.md"}
	for _, file := range files {
		if !slices.Contains(expectedFilenames, file.Name()) {
			t.Errorf("Unexpected snapshot filename: %s", file.Name())
		}
	}

	// スナップショットの内容を確認
	for _, file := range files {
		content, err := os.ReadFile(filepath.Join(snapshotDir, file.Name()))
		if err != nil {
			t.Fatalf("Failed to read snapshot file: %v", err)
		}

		if string(content) != "# Test Article\n\nThis is a test article." {
			t.Errorf("Unexpected content: %s", string(content))
		}
	}
}
