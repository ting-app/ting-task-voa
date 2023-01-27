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

	words, err := parseWords(bodyNodes)

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
		Words:        words,
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

func parseWords(nodes []*html.Node) ([]Word, error) {
	var words []Word

	for i := 0; i < len(nodes); i++ {
		text := goquery.NewDocumentFromNode(nodes[i]).Text()

		if text == "Words in This Story" {
			startIndex := i + 1

			for j := startIndex; j < len(nodes); j++ {
				wordNode := goquery.NewDocumentFromNode(nodes[j])

				if strings.HasPrefix(wordNode.Text(), "___") {
					break
				}

				word := parseWord(wordNode)

				if word != nil {
					words = append(words, *word)
				}
			}

			break
		}
	}

	if len(words) == 0 {
		return nil, errors.New("failed to parse words")
	}

	return words, nil
}

func parseWord(document *goquery.Document) *Word {
	var words []string

	// It might be a phrase
	document.Find("strong").Each(func(i int, selection *goquery.Selection) {
		wordNode := goquery.NewDocumentFromNode(selection.Nodes[0])

		words = append(words, wordNode.Text())
	})

	word := strings.TrimSpace(strings.Join(words, " "))

	if word == "" {
		return nil
	}

	all := document.Text()
	rest := strings.Split(all, word)[1]
	partOfSpeech := strings.TrimSpace(rest[0 : strings.Index(rest, ".")+1])

	if strings.HasPrefix(partOfSpeech, "–") {
		partOfSpeech = strings.Split(partOfSpeech, "–")[1]
	}

	definition := strings.TrimSpace(rest[strings.Index(rest, ".")+1:])

	return &Word{
		Word:         word,
		PartOfSpeech: partOfSpeech,
		Definition:   definition,
	}
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
