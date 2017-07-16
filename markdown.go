package main

import "github.com/russross/blackfriday"

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
