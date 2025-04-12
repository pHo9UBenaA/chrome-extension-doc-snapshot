package storage

import (
	"fmt"
	"os"

	"path/filepath"
	"strings"
)

const (
	// SnapshotDirはスナップショットディレクトリ
	SnapshotDir = "__snapshot__"

	// SnapshotFileExtensionはスナップショットファイルの拡張子
	SnapshotFileExtension = ".md"
)

func getSnapshotDir() string {
	if dir := os.Getenv("SNAPSHOT_DIR"); dir != "" {
		return filepath.Join(dir, SnapshotDir)
	}
	return SnapshotDir
}

// EnsureSnapshotDirはスナップショットディレクトリが存在することを保証します
func EnsureSnapshotDir() error {
	if err := os.MkdirAll(getSnapshotDir(), 0755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %v", err)
	}
	return nil
}

// TakeSnapshotは指定された内容をスナップショットファイルとして保存します
func TakeSnapshot(href string, markdown string) error {
	// ディレクトリの存在を確認
	if err := EnsureSnapshotDir(); err != nil {
		return err
	}

	// ファイル名を生成（最後の部分を使用）
	fileName := href
	if idx := strings.LastIndex(href, "/"); idx != -1 {
		fileName = href[idx+1:]
	}
	fileName = fileName + SnapshotFileExtension

	// ファイルパスを生成
	filePath := filepath.Join(getSnapshotDir(), fileName)

	// ファイルに書き込み
	if err := os.WriteFile(filePath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %v", err)
	}

	return nil
}
