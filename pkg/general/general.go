package general

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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

func DecryptAES(cipherText string, key string) (string, error) {
	keyBytes := []byte(key)

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce := data[:nonceSize]
	cipherBytes := data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
