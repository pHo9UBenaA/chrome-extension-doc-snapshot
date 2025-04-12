package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_TakeSnapshot(t *testing.T) {
	// Arrange
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	os.Setenv("SNAPSHOT_DIR", tempDir)

	// Assert
	// テストデータ
	href := "/test/article"
	markdown := "# Test Article\n\nThis is a test article."

	// Act
	err := TakeSnapshot(href, markdown)

	// Assert
	if err != nil {
		t.Fatalf("TakeSnapshot failed: %v", err)
	}

	// ファイルが作成されたことを確認
	filePath := filepath.Join(GetSnapshotDirPath(), "article.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("Snapshot file was not created")
	}

	// ファイルの内容を確認
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read snapshot file: %v", err)
	}

	if string(content) != markdown {
		t.Errorf("Expected content '%s', got '%s'", markdown, string(content))
	}
}
