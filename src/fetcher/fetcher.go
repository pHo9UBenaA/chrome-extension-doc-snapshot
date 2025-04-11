package fetcher

import (
	"net/http"
	"time"

	"golang.org/x/net/html"
)

// FetchHTMLは指定されたURLからHTMLを取得し、パースして返します
func FetchHTML(url string, timeout time.Duration) (*html.Node, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return html.Parse(res.Body)
}
