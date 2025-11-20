package general

import (
	"bytes"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func ParseTemplateEmailToPlainText(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		if n.Type == html.ElementNode {
			switch n.Data {
			case "br", "p", "div", "hr":
				buf.WriteString("\n")
			case "a":
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						buf.WriteString(" [" + attr.Val + "]")
						break
					}
				}
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	output := buf.String()
	output = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(output, "\n\n")
	return strings.TrimSpace(output)
}
