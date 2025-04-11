package crawler

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/converter"
	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/fetcher"
	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/parser"
	"github.com/pHo9UBenaA/chrome-extension-doc-snapshot/src/storage"
)

// Crawlerはクローリングのコアロジックを提供します
type Crawler struct {
	storage storage.Storage
}

// NewCrawlerは新しいCrawlerインスタンスを作成します
func NewCrawler(storage storage.Storage) *Crawler {
	return &Crawler{
		storage: storage,
	}
}

// FetchAndSnapshotArticleは指定されたURLから記事を取得し、スナップショットとして保存します
func (c *Crawler) FetchAndSnapshotArticle(url string) error {
	// HTMLを取得
	doc, err := fetcher.FetchHTML(url, 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to fetch HTML: %v", err)
	}

	// 記事を抽出
	articleNode, err := parser.ExtractArticle(doc)
	if err != nil {
		return fmt.Errorf("failed to extract article: %v", err)
	}

	// Markdownに変換
	markdown, err := converter.ConvertNodeToMarkdown(articleNode)
	if err != nil {
		return fmt.Errorf("failed to convert to markdown: %v", err)
	}

	// URLからパス部分を抽出し、最後の部分だけを取得
	path := url
	if idx := strings.LastIndex(path, "/"); idx != -1 {
		path = path[idx+1:]
	}

	// スナップショットとして保存
	if err := c.storage.TakeSnapshot(path, markdown); err != nil {
		return fmt.Errorf("failed to take snapshot: %v", err)
	}

	return nil
}

// FetchAPIReferenceはAPIリファレンスページを取得します
func (c *Crawler) FetchAPIReference(baseURL string) (*html.Node, error) {
	return fetcher.FetchHTML(baseURL+"/docs/extensions/reference/api", 30*time.Second)
}
