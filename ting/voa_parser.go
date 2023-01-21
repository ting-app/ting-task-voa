package ting

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"time"
)

type Voa struct {
	Title          string
	Description    string
	Url            string
	PublishedAtUtc time.Time
	Body           string
	BodyWithHtml   string
	ImageUrl       string
}

func parseVoa(url string) (*Voa, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status code error: %d %s", response.StatusCode, response.Status))
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return nil, err
	}

	contentNodes := doc.Find("#article-content").Nodes

	if len(contentNodes) == 0 {
		return nil, errors.New("article content not found")
	}

	contentNode := goquery.NewDocumentFromNode(contentNodes[0]).Children().Nodes[0]
	bodyWithHtml, body, err := parseContent(goquery.NewDocumentFromNode(contentNode).Children().Nodes)

	if err != nil {
		return nil, err
	}

	voa := &Voa{
		Url:          url,
		Body:         strings.Join(body, "\n"),
		BodyWithHtml: strings.Join(bodyWithHtml, "\n"),
	}

	return voa, nil
}

func parseContent(nodes []*html.Node) ([]string, []string, error) {
	var contentNodes []*html.Node

	// The first node is audio node
	for i := 1; i < len(nodes); i++ {
		node := nodes[i]

		if node.Data == "p" {
			contentNodes = append(contentNodes, node)
		} else if node.Data == "h2" {
			break
		}
	}

	var bodyWithHtml []string
	var body []string

	for _, node := range contentNodes {
		wrapperNode := goquery.NewDocumentFromNode(node)
		bodyHtml, err := wrapperNode.Html()

		if err != nil {
			return nil, nil, errors.New("failed to parse html from node")
		}

		content := wrapperNode.Text()

		// End of content
		if strings.HasPrefix(content, "___") {
			break
		}

		bodyWithHtml = append(bodyWithHtml, fmt.Sprintf("<%s>%s<%s>", node.Data, bodyHtml, node.Data))
		body = append(body, content)
	}

	return bodyWithHtml, body, nil
}
