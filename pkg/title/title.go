/* Based on: https://siongui.github.io/2016/05/10/go-get-html-title-via-net-html/ */
package title

import (
	"errors"
	"fmt"
	"io"

	"golang.org/x/net/html"
)

func traverse(n *html.Node) (string, error) {
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data, nil
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, err := traverse(c)
		if err == nil {
			return result, err
		}
	}

	return "", errors.New("Could not find title in html")
}

func GetHtmlTitle(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", fmt.Errorf("Error parsing HTML: %s", err)
	}

	return traverse(doc)
}
