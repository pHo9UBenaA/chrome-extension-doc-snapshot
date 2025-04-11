package storage

import (
	"fmt"
	"os"

	// "path"
	"path/filepath"
	"strings"
)

// Storageはファイルシステム操作のインターフェース
type Storage interface {
	EnsureSnapshotDir() error
	TakeSnapshot(href string, markdown string) error
}

// FileStorageは実際のファイルシステムを使用する実装
type FileStorage struct {
	SnapshotDir           string
	SnapshotFileExtension string
}

// NewFileStorageは新しいFileStorageインスタンスを作成します
func NewFileStorage() *FileStorage {
	return &FileStorage{
		SnapshotDir:           "__snapshot__",
		SnapshotFileExtension: ".md",
	}
}

// EnsureSnapshotDirはスナップショットディレクトリが存在することを保証します
func (s *FileStorage) EnsureSnapshotDir() error {
	if err := os.MkdirAll(s.SnapshotDir, 0755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %v", err)
	}
	return nil
}

// TakeSnapshotは指定された内容をスナップショットファイルとして保存します
func (s *FileStorage) TakeSnapshot(href string, markdown string) error {
	// ディレクトリの存在を確認
	if err := s.EnsureSnapshotDir(); err != nil {
		return err
	}

	// ファイル名を生成（最後の部分を使用）
	fileName := href
	if idx := strings.LastIndex(href, "/"); idx != -1 {
		fileName = href[idx+1:]
	}
	fileName = fileName + s.SnapshotFileExtension

	// ファイルパスを生成
	filePath := filepath.Join(s.SnapshotDir, fileName)

	// ファイルに書き込み
	if err := os.WriteFile(filePath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %v", err)
	}

	return nil
}

// // MockStorageはテスト用のモック実装
// type MockStorage struct {
// 	Snapshots map[string]string
// }

// // NewMockStorageは新しいMockStorageインスタンスを作成します
// func NewMockStorage() *MockStorage {
// 	return &MockStorage{
// 		Snapshots: make(map[string]string),
// 	}
// }

// // EnsureSnapshotDirはモック実装
// func (s *MockStorage) EnsureSnapshotDir() error {
// 	return nil
// }

// // TakeSnapshotはモック実装
// func (s *MockStorage) TakeSnapshot(href string, markdown string) error {
// 	s.Snapshots[path.Base(href)] = markdown
// 	return nil
// }
