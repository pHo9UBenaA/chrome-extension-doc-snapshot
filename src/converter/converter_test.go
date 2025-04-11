package converter

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestConvertNodeToMarkdown(t *testing.T) {
	// Arrange
	// テスト用のHTMLノードを作成
	doc := &html.Node{
		Type: html.ElementNode,
		Data: "div",
		FirstChild: &html.Node{
			Type: html.ElementNode,
			Data: "h1",
			FirstChild: &html.Node{
				Type: html.TextNode,
				Data: "Test Title",
			},
		},
	}

	// Act
	markdown, err := ConvertNodeToMarkdown(doc)

	// Assert
	if err != nil {
		t.Fatalf("ConvertNodeToMarkdown failed: %v", err)
	}

	// 期待されるMarkdownの形式を確認
	expected := "# Test Title"
	if !strings.Contains(markdown, expected) {
		t.Errorf("Expected markdown to contain '%s', got '%s'", expected, markdown)
	}
}

func TestRenderNode(t *testing.T) {
	// Arrange
	node := &html.Node{
		Type: html.ElementNode,
		Data: "p",
		FirstChild: &html.Node{
			Type: html.TextNode,
			Data: "Test paragraph",
		},
	}

	// Act
	htmlContent := RenderNode(node)

	// Assert
	expected := "<p>Test paragraph</p>"
	if htmlContent != expected {
		t.Errorf("Expected HTML '%s', got '%s'", expected, htmlContent)
	}
}
