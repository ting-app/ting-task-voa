package ting

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

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

	articleContentNodes := doc.Find("#article-content").Nodes

	if len(articleContentNodes) == 0 {
		return nil, errors.New("article content not found")
	}

	articleContentNode := goquery.NewDocumentFromNode(articleContentNodes[0]).Children().Nodes[0]
	bodyNodes := goquery.NewDocumentFromNode(articleContentNode).Children().Nodes
	audioUrl, err := parseAudioUrl(bodyNodes[0])

	if err != nil {
		return nil, err
	}

	bodyWithHtml, body, err := parseContent(bodyNodes)

	if err != nil {
		return nil, err
	}

	voa := &Voa{
		Url:          url,
		Body:         strings.Join(body, "\n"),
		BodyWithHtml: strings.Join(bodyWithHtml, "\n"),
		AudioUrl:     audioUrl,
	}

	return voa, nil
}

func parseAudioUrl(node *html.Node) (string, error) {
	downloadNodes := goquery.NewDocumentFromNode(node).Find(".media-download li.subitem a").Nodes

	if len(downloadNodes) == 0 {
		return "", errors.New("no download nodes found")
	}

	attributes := downloadNodes[0].Attr

	for _, attribute := range attributes {
		if attribute.Key == "href" {
			url := attribute.Val

			if strings.Contains(url, "?") {
				url = strings.Split(url, "?")[0]
			}

			return url, nil
		}
	}

	return "", errors.New("no audio url found")
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
