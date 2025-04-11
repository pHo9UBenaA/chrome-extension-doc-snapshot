package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/storage/mock"
)

func TestFileStorage(t *testing.T) {
	// Arrange
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	storage := &FileStorage{
		SnapshotDir:           filepath.Join(tempDir, "__snapshot__"),
		SnapshotFileExtension: ".md",
	}

	// Act
	err := storage.EnsureSnapshotDir()

	// Assert
	if err != nil {
		t.Fatalf("EnsureSnapshotDir failed: %v", err)
	}

	// ディレクトリが存在することを確認
	if _, err := os.Stat(storage.SnapshotDir); os.IsNotExist(err) {
		t.Fatal("Snapshot directory was not created")
	}

	// テストデータ
	href := "/test/article"
	markdown := "# Test Article\n\nThis is a test article."

	// Act
	err = storage.TakeSnapshot(href, markdown)

	// Assert
	if err != nil {
		t.Fatalf("TakeSnapshot failed: %v", err)
	}

	// ファイルが作成されたことを確認
	filePath := filepath.Join(storage.SnapshotDir, "article.md")
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

func TestMockStorage(t *testing.T) {
	// Arrange
	storage := mock.NewMockStorage()

	// テストデータ
	href := "/test/article"
	markdown := "# Test Article\n\nThis is a test article."

	// Act
	err := storage.TakeSnapshot(href, markdown)

	// Assert
	if err != nil {
		t.Fatalf("TakeSnapshot failed: %v", err)
	}

	// スナップショットが保存されたことを確認
	savedMarkdown, exists := storage.Snapshots["article"]
	if !exists {
		t.Fatal("Snapshot was not saved")
	}

	if savedMarkdown != markdown {
		t.Errorf("Expected content '%s', got '%s'", markdown, savedMarkdown)
	}
}
