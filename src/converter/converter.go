package converter

import (
	"bytes"

	"golang.org/x/net/html"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// renderNodeはHTMLノードを文字列に変換します
func renderNode(n *html.Node) string {
	buf := new(bytes.Buffer)
	html.Render(buf, n)
	return buf.String()
}

// NodeToMarkdownはHTMLノードをMarkdownに変換します
func NodeToMarkdown(n *html.Node) (string, error) {
	htmlContent := renderNode(n)
	return htmltomarkdown.ConvertString(htmlContent)
}
