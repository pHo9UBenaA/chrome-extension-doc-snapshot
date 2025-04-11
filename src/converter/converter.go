package converter

import (
	"bytes"

	"golang.org/x/net/html"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// ConvertNodeToMarkdownはHTMLノードをMarkdownに変換します
func ConvertNodeToMarkdown(n *html.Node) (string, error) {
	htmlContent := RenderNode(n)
	return htmltomarkdown.ConvertString(htmlContent)
}

// RenderNodeはHTMLノードを文字列に変換します
func RenderNode(n *html.Node) string {
	buf := new(bytes.Buffer)
	html.Render(buf, n)
	return buf.String()
}
