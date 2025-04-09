package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
)

func setupTempDir(_ *testing.M) (string, string) {
	// `t.tempDir`でテスト毎に一時的な領域を作るのが面倒なためTestMainで設定している
	// 多分これのせいでテストがアイソレートされていないっぽい

	origWd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get working directory: %v\n", err)
		os.Exit(1)
	}

	tempWd, err := os.MkdirTemp("", "test")
	if err != nil {
		fmt.Printf("failed to create temporary directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chdir(tempWd); err != nil {
		fmt.Printf("failed to change working directory: %v\n", err)
		os.Exit(1)
	}

	return origWd, tempWd
}

func shutdown(origWd, tempWd string) {
	os.RemoveAll(tempWd)
	os.Chdir(origWd)
}

func TestMain(m *testing.M) {
	origWd, tempWd := setupTempDir(m)

	// Go1.5からos.Exitは不要になったよう
	// https://go.dev/doc/go1.15#testing
	m.Run()

	shutdown(origWd, tempWd)
}

func Test_EnsureSnapshotDir(t *testing.T) {
	if err := ensureSnapshotDir(); err != nil {
		t.Fatalf("%v", err)
	}

	info, err := os.Stat(SnapshotDir)
	if err != nil {
		t.Fatalf("failed to check snapshot directory: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("expected %s to be a directory", SnapshotDir)
	}
}

func Test_TakeSnapshot(t *testing.T) {
	fileName := "test"
	// HTMLコンテンツをテスト
	content := "<article><h1>Test</h1><p>Hello, world!</p></article>"
	expected := "# Test\n\nHello, world!"

	if err := takeSnapshot(fileName, content); err != nil {
		t.Fatalf("%v", err)
	}

	snapshotFile := path.Join(SnapshotDir, fileName+SnapshotFileExtension)
	data, err := os.ReadFile(snapshotFile)

	if err != nil {
		t.Fatalf("failed to read snapshot file: %v", err)
	}

	if string(data) != expected {
		t.Errorf("expected content %q, got %q", expected, string(data))
	}
}

func Test_ScrapeAPIReference(t *testing.T) {
	links, err := fetchAPIReference()

	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(links) == 0 {
		t.Errorf("expected at least one link, got 0")
	}
}

func Test_ScrapeArticle(t *testing.T) {
	href := "/docs/extensions/reference/api/tabs"
	content, err := fetchArticle(href)

	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(content) == 0 {
		t.Errorf("expected content, got empty string")
	}

	// HTMLの基本構造をチェック
	if !strings.Contains(content, "<article") || !strings.Contains(content, "</article>") {
		t.Errorf("expected HTML article tags, got %q", content)
	}
}
