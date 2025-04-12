package storage

import (
	"fmt"
	"os"

	"path/filepath"
	"strings"
)

const (
	// SnapshotDirはスナップショットディレクトリ
	snapshotDir = "__snapshot__"
	// snapshotFileExtensionはスナップショットファイルの拡張子
	snapshotFileExtension = ".md"
)

func GetSnapshotDirPath() string {
	if dir := os.Getenv("SNAPSHOT_DIR"); dir != "" {
		return filepath.Join(dir, snapshotDir)
	}
	return snapshotDir
}

// EnsureSnapshotDirはスナップショットディレクトリが存在することを保証します
func ensureSnapshotDir() error {
	if err := os.MkdirAll(GetSnapshotDirPath(), 0755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %v", err)
	}
	return nil
}

// TakeSnapshotは指定された内容をスナップショットファイルとして保存します
func TakeSnapshot(href string, markdown string) error {
	// ディレクトリの存在を確認
	if err := ensureSnapshotDir(); err != nil {
		return err
	}

	// ファイル名を生成（最後の部分を使用）
	fileName := href
	if idx := strings.LastIndex(href, "/"); idx != -1 {
		fileName = href[idx+1:]
	}
	fileName = fileName + snapshotFileExtension

	// ファイルパスを生成
	filePath := filepath.Join(GetSnapshotDirPath(), fileName)

	// ファイルに書き込み
	if err := os.WriteFile(filePath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %v", err)
	}

	return nil
}
