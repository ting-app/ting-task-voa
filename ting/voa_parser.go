package ting

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func parseVoa(url string) (*Voa, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status code error: %d %s, url=%s", response.StatusCode, response.Status, url))
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return nil, err
	}

	articleContentNodes := doc.Find("#article-content").Nodes

	if len(articleContentNodes) == 0 {
		return nil, errors.New(fmt.Sprintf("article content not found, url=%s", url))
	}

	articleContentNode := goquery.NewDocumentFromNode(articleContentNodes[0]).Children().Nodes[0]
	bodyNodes := parseBodyNodes(articleContentNode)
	audioUrl, err := parseAudioUrl(bodyNodes[0])

	if err != nil {
		return nil, err
	}

	words, err := parseWords(url, bodyNodes)

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

func parseWords(url string, nodes []*html.Node) ([]Word, error) {
	var words []Word

	for i := 0; i < len(nodes); i++ {
		text := goquery.NewDocumentFromNode(nodes[i]).Text()

		if strings.Contains(strings.TrimSpace(text), "Words in This Story") {
			startIndex := i + 1

			for j := startIndex; j < len(nodes); j++ {
				wordNode := goquery.NewDocumentFromNode(nodes[j])

				if wordNode.Nodes[0].Data == "div" {
					for k := 0; k < wordNode.Children().Length(); k++ {
						childNode := goquery.NewDocumentFromNode(wordNode.Children().Nodes[k])
						word := addWord(url, childNode)

						if word != nil {
							words = append(words, *word)
						}
					}
				} else {
					word := addWord(url, wordNode)

					if word != nil {
						words = append(words, *word)
					}
				}
			}

			break
		}
	}

	if len(words) == 0 {
		return nil, errors.New(fmt.Sprintf("failed to parse words, url=%s", url))
	}

	return words, nil
}

func addWord(url string, wordNode *goquery.Document) *Word {
	wordText := wordNode.Text()

	if strings.HasPrefix(wordText, "___") {
		return nil
	}

	word := parseWord(url, wordNode)

	return word
}

func parseWord(url string, document *goquery.Document) *Word {
	var words []string

	// It might be a phrase
	document.Find("strong").Each(func(i int, selection *goquery.Selection) {
		wordNode := goquery.NewDocumentFromNode(selection.Nodes[0])
		wordText := strings.TrimSpace(wordNode.Text())

		matched, _ := regexp.MatchString("\\w+", wordText)

		if matched {
			words = append(words, wordText)
		} else {
			log.Printf("Unmatched word %s", wordText)
		}
	})

	word := strings.TrimSpace(strings.Join(words, " "))

	if word == "" {
		return nil
	}

	all := document.Text()
	rests := strings.Split(all, word)

	if len(rests) == 1 {
		log.Println("h")
	}

	rest := rests[1]
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

func parseBodyNodes(articleContentNode *html.Node) []*html.Node {
	children := goquery.NewDocumentFromNode(articleContentNode).Children()

	if children.Length() <= 2 {
		return goquery.NewDocumentFromNode(children.Nodes[0]).Children().Nodes
	} else {
		return children.Nodes
	}
}
