package mock

import (
	"strings"
)

// MockStorage はテスト用のモックストレージ実装
type MockStorage struct {
	Snapshots map[string]string
}

// NewMockStorage は新しいMockStorageインスタンスを作成
func NewMockStorage() *MockStorage {
	return &MockStorage{
		Snapshots: make(map[string]string),
	}
}

// TakeSnapshot はスナップショットをメモリに保存
func (s *MockStorage) TakeSnapshot(href string, markdown string) error {
	// ファイル名を生成（最後の部分を使用）
	fileName := href
	if idx := strings.LastIndex(href, "/"); idx != -1 {
		fileName = href[idx+1:]
	}
	s.Snapshots[fileName] = markdown
	return nil
}

// EnsureSnapshotDir はモック用の空実装
func (s *MockStorage) EnsureSnapshotDir() error {
	return nil
}
