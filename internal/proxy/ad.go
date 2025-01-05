package proxy

import (
	"bytes"
	"log"

	"golang.org/x/net/html"
)

func removeAds(body []byte) ([]byte, error) {
	log.Printf("Modifying HTML to remove iframe and ins tags...")

	rootNode, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "iframe":
				log.Print("Gottem iframe")
				n.Parent.RemoveChild(n)
			case "ins":
				log.Print("ins, gottem!")
				n.Parent.RemoveChild(n)
			}
		}

		for c := n.FirstChild; c != nil; {
			next := c.NextSibling
			traverse(c)
			c = next
		}
	}

	traverse(rootNode)

	var buf bytes.Buffer
	err = html.Render(&buf, rootNode)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
