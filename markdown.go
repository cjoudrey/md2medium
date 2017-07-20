package main

import "regexp"
import "github.com/russross/blackfriday"
import "gopkg.in/yaml.v2"

type MarkdownMetadata struct {
	Tags []string `yaml:"tags"`
}

const extensions = blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_SPACE_HEADERS |
	blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK

func MarkdownToHTML(markdown []byte) []byte {
	const htmlFlags = 0
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	html := blackfriday.Markdown(markdown, renderer, extensions)
	return html
}

func ExtractMarkdownMetadata(markdown []byte) (*MarkdownMetadata, []byte, error) {
	var metadata MarkdownMetadata

	splitter := regexp.MustCompile("(?s)^---\n(.*)---\n(.*)")
	matches := splitter.FindAllSubmatch(markdown, -1)

	if len(matches) == 0 {
		return &metadata, markdown, nil
	}

	yamlData := matches[0][1]
	if err := yaml.Unmarshal(yamlData, &metadata); err != nil {
		return nil, nil, err
	}

	strippedMarkdown := matches[0][2]
	return &metadata, strippedMarkdown, nil

}
