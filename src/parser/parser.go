package parser

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// removeUnwantedElementsは不要な要素を削除します
func removeUnwantedElements(n *html.Node) {
	var next *html.Node

	for c := n.FirstChild; c != nil; c = next {
		next = c.NextSibling

		if c.Type != html.ElementNode {
			continue
		}

		switch {
		// パンくずリスト
		case c.Data == "div" && hasClass(c, "devsite-article-meta"):
			if c.Parent != nil {
				c.Parent.RemoveChild(c)
			}
		// ヘッダーのツールチップ
		case c.Data == "h1" && hasClass(c, "devsite-page-title"):
			for inner := c.FirstChild; inner != nil; inner = inner.NextSibling {
				if inner.Data == "div" && inner.Parent != nil {
					inner.Parent.RemoveChild(inner)
				}
			}
		}
	}
}

// hasClassはノードが指定されたクラスを持っているかどうかを返します
func hasClass(n *html.Node, className string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			return strings.Contains(attr.Val, className)
		}
	}
	return false
}

// findHrefInAnchorはアンカータグからhref属性の値を取得します
func findHrefInAnchor(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				return attr.Val
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if href := findHrefInAnchor(c); href != "" {
			return href
		}
	}
	return ""
}

// findArticleはHTMLドキュメントからarticle要素を探して返します
func findArticle(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "article" {
		removeUnwantedElements(n)
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findArticle(c); result != nil {
			return result
		}
	}
	return nil
}

// findAPILinksはHTMLドキュメントからAPIリンクを抽出します
func findAPILinks(n *html.Node) []string {
	var apiLinks []string

	if n.Type == html.ElementNode && n.Data == "dt" {
		if href := findHrefInAnchor(n); href != "" {
			apiLinks = append(apiLinks, href)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		apiLinks = append(apiLinks, findAPILinks(c)...)
	}

	return apiLinks
}

// ExtractArticleはHTMLドキュメントから記事を抽出します
func ExtractArticle(doc *html.Node) (*html.Node, error) {
	if article := findArticle(doc); article != nil {
		return article, nil
	}
	return nil, fmt.Errorf("article not found")
}

// ExtractAPILinksはHTMLドキュメントからAPIリンクを抽出します
func ExtractAPILinks(doc *html.Node) ([]string, error) {
	if apiLinks := findAPILinks(doc); len(apiLinks) > 0 {
		return apiLinks, nil
	}
	return nil, fmt.Errorf("API links not found")
}
