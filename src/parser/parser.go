package parser

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// FindArticleはHTMLドキュメントからarticle要素を探して返します
func FindArticle(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "article" {
		RemoveUnwantedElements(n)
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := FindArticle(c); result != nil {
			return result
		}
	}
	return nil
}

// ExtractArticleはHTMLドキュメントから記事を抽出します
func ExtractArticle(doc *html.Node) (*html.Node, error) {
	if article := FindArticle(doc); article != nil {
		return article, nil
	}
	return nil, fmt.Errorf("article not found")
}

// RemoveUnwantedElementsは不要な要素を削除します
func RemoveUnwantedElements(n *html.Node) {
	var next *html.Node

	for c := n.FirstChild; c != nil; c = next {
		next = c.NextSibling

		if c.Type != html.ElementNode {
			continue
		}

		switch {
		// パンくずリスト
		case c.Data == "div" && HasClass(c, "devsite-article-meta"):
			if c.Parent != nil {
				c.Parent.RemoveChild(c)
			}
		// ヘッダーのツールチップ
		case c.Data == "h1" && HasClass(c, "devsite-page-title"):
			for inner := c.FirstChild; inner != nil; inner = inner.NextSibling {
				if inner.Data == "div" && inner.Parent != nil {
					inner.Parent.RemoveChild(inner)
				}
			}
		}
	}
}

// HasClassはノードが指定されたクラスを持っているかどうかを返します
func HasClass(n *html.Node, className string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			return strings.Contains(attr.Val, className)
		}
	}
	return false
}

// FindHrefInAnchorはアンカータグからhref属性の値を取得します
func FindHrefInAnchor(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				return attr.Val
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if href := FindHrefInAnchor(c); href != "" {
			return href
		}
	}
	return ""
}

// ExtractAPILinksはHTMLドキュメントからAPIリンクを抽出します
func ExtractAPILinks(doc *html.Node) []string {
	var apiLinks []string

	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "dt" {
			if href := FindHrefInAnchor(n); href != "" {
				apiLinks = append(apiLinks, href)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}

	findLinks(doc)
	return apiLinks
}
