package parser

import (
	"testing"

	"golang.org/x/net/html"
)

func Test_findArticle(t *testing.T) {
	// Arrange
	// テスト用のHTMLノードを作成
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "body",
			FirstChild: &html.Node{
				Type: html.ElementNode,
				Data: "article",
				Attr: []html.Attribute{
					{Key: "class", Val: "devsite-article"},
				},
			},
		},
	}

	// Act
	article := findArticle(doc)

	// Assert
	if article == nil {
		t.Fatal("Expected to find article node")
	}

	if article.Data != "article" {
		t.Errorf("Expected node type 'article', got '%s'", article.Data)
	}
}

func Test_removeUnwantedElements(t *testing.T) {
	// Arrange
	// 不要な要素を含むHTMLノードを作成
	parent := &html.Node{
		Type: html.ElementNode,
		Data: "article",
	}
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			{Key: "class", Val: "devsite-article-meta"},
		},
		Parent: parent,
	}
	parent.FirstChild = doc

	// Act
	removeUnwantedElements(parent)

	// Assert
	// 不要な要素が削除されていることを確認
	if parent.FirstChild != nil {
		t.Error("Expected unwanted elements to be removed")
	}
}

func Test_hasClass(t *testing.T) {
	// Arrange
	node := &html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			{Key: "class", Val: "test-class another-class"},
		},
	}

	// Act & Assert
	// クラスが存在する場合
	if !hasClass(node, "test-class") {
		t.Error("Expected to find 'test-class'")
	}

	// クラスが存在しない場合
	if hasClass(node, "non-existent") {
		t.Error("Expected not to find 'non-existent'")
	}
}

func Test_findHrefInAnchor(t *testing.T) {
	// Arrange
	node := &html.Node{
		Type: html.ElementNode,
		Data: "dt",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "a",
			Attr: []html.Attribute{
				{Key: "href", Val: "/test/link"},
			},
		},
	}

	// Act
	href := findHrefInAnchor(node)

	// Assert
	if href != "/test/link" {
		t.Errorf("Expected href '/test/link', got '%s'", href)
	}
}

func Test_ExtractArticle(t *testing.T) {
	// Arrange
	// テスト用のHTMLノードを作成
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "body",
			FirstChild: &html.Node{
				Type: html.ElementNode,
				Data: "article",
			},
		},
	}

	// Act
	article, err := ExtractArticle(doc)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if article == nil {
		t.Fatal("Expected to find article node")
	}
}

func Test_ExtractAPILinks(t *testing.T) {
	// Arrange
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "body",
			FirstChild: &html.Node{
				Type: html.ElementNode,
				Data: "dl",
				FirstChild: &html.Node{
					Type: html.ElementNode,
					Data: "dt",
					FirstChild: &html.Node{
						Type: html.ElementNode,
						Data: "a",
						Attr: []html.Attribute{
							{Key: "href", Val: "/api1"},
						},
					},
				},
			},
		},
	}

	// Act
	links, err := ExtractAPILinks(doc)
	if err != nil {
		t.Fatalf("ExtractAPILinks failed: %v", err)
	}

	// Assert
	if len(links) != 1 {
		t.Fatalf("Expected 1 link, got %d", len(links))
	}

	if links[0] != "/api1" {
		t.Errorf("Expected link '/api1', got '%s'", links[0])
	}
}
