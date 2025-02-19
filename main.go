package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	SnapshotDir = "__snapshot__"
	BaseURL     = "https://developer.chrome.com/"
)

func main() {
	log.Println("start")

	// snapshotディレクトリの作成（存在しない場合）
	if err := os.MkdirAll(SnapshotDir, 0755); err != nil {
		log.Fatalf("failed to create snapshot directory: %v", err)
	}

	hrefList, err := getReferenceList()
	if err != nil {
		log.Fatalf("failed to get reference list: %v", err)
	}

	log.Println("Done: getReferenceList")

	for _, href := range hrefList {
		log.Printf("start: getArticle %s\n", href)
		if err := getArticle(href); err != nil {
			log.Printf("failed to get article (%s): %v", href, err)
			continue
		}
		log.Printf("Done: getArticle %s\n", href)

		time.Sleep(1 * time.Second)
	}

	log.Println("Done: main")
}

func getReferenceList() ([]string, error) {
	var refList []string

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(2),
	)

	c.OnHTML("dl", func(e *colly.HTMLElement) {
		e.ForEach("dt", func(_ int, el *colly.HTMLElement) {
			href := el.ChildAttr("a", "href")
			if href != "" {
				refList = append(refList, href)
			}
		})
	})

	visitURL := BaseURL + "docs/extensions/reference/api"
	if err := c.Visit(visitURL); err != nil {
		return nil, fmt.Errorf("visit error: %w", err)
	}

	// インデックスファイルの保存
	indexPath := path.Join(SnapshotDir, "index.txt")
	content := strings.Join(refList, "\n")
	if err := os.WriteFile(indexPath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write index file: %w", err)
	}

	return refList, nil
}

func getArticle(href string) error {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(2),
	)

	var articleContent string
	c.OnHTML("article", func(e *colly.HTMLElement) {
		articleContent = e.Text
	})

	if err := c.Visit(BaseURL + href); err != nil {
		return fmt.Errorf("failed to visit article: %w", err)
	}

	baseName := path.Base(href)
	filePath := path.Join(SnapshotDir, baseName+".txt")
	if err := os.WriteFile(filePath, []byte(articleContent), 0644); err != nil {
		return fmt.Errorf("failed to write article file: %w", err)
	}

	return nil
}