package main

import "github.com/PuerkitoBio/goquery"
import "golang.org/x/net/html"
import "strings"
import "path"
import "github.com/google/go-github/github"
import "github.com/Medium/medium-sdk-go"
import "context"
import "fmt"
import "mime"

type HtmlFormatter struct {
	path   string
	doc    *goquery.Document
	logger *PrettyLogger
	ctx    context.Context
}

func NewHtmlFormatter(filePath string, html string, logger *PrettyLogger, ctx context.Context) (*HtmlFormatter, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	htmlFormatter := HtmlFormatter{
		path:   path.Dir(filePath),
		doc:    doc,
		logger: logger,
		ctx:    ctx,
	}

	return &htmlFormatter, nil
}

func (htmlFormatter *HtmlFormatter) Html() (string, error) {
	return htmlFormatter.doc.Html()
}

func (htmlFormatter *HtmlFormatter) ReplaceCodeBlocks(githubClient *github.Client) error {
	codeBlocks := htmlFormatter.doc.Find("pre code[class^='language-']")

	for i := 0; i < codeBlocks.Length(); i++ {
		codeBlock := codeBlocks.Eq(i)
		className, _ := codeBlock.Attr("class")

		language := className[9:]
		code := codeBlock.Text()

		htmlFormatter.logger.Loading()
		htmlFormatter.logger.Info("Creating secret gist")
		public := false
		gist, _, err := githubClient.Gists.Create(htmlFormatter.ctx, &github.Gist{
			Public: &public,
			Files: map[github.GistFilename]github.GistFile{
				github.GistFilename("snippet." + language): github.GistFile{
					Content: &code,
				},
			},
		})
		htmlFormatter.logger.Done()
		if err != nil {
			return err
		}

		gistLink := html.Node{
			Type: html.ElementNode,
			Data: "a",
			Attr: []html.Attribute{
				html.Attribute{Key: "href", Val: *gist.HTMLURL},
			},
		}

		gistLink.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: *gist.HTMLURL,
		})

		codeBlock.Parent().ReplaceWithNodes(&gistLink)
	}

	return nil
}

func (htmlFormatter *HtmlFormatter) ReplaceImages(mediumClient *medium.Medium) error {
	err := htmlFormatter.UploadLocalImages(mediumClient)
	if err != nil {
		return err
	}

	htmlFormatter.DisplayAltAsFigcaption()

	return nil
}

func (htmlFormatter *HtmlFormatter) UploadLocalImages(mediumClient *medium.Medium) error {
	localImages := htmlFormatter.doc.Find("img:not([src^='http://']):not([src^='https://'])")

	for i := 0; i < localImages.Length(); i++ {
		localImage := localImages.Eq(i)

		src, _ := localImage.Attr("src")
		absSrc := path.Join(htmlFormatter.path, src)

		htmlFormatter.logger.Loading()
		htmlFormatter.logger.Info(fmt.Sprintf("Uploading image: %s", absSrc))
		image, err := mediumClient.UploadImage(medium.UploadOptions{
			FilePath:    absSrc,
			ContentType: mime.TypeByExtension(path.Ext(src)),
		})
		htmlFormatter.logger.Done()
		if err != nil {
			return err
		}

		localImage.SetAttr("src", image.URL)
	}

	return nil
}

func (htmlFormatter *HtmlFormatter) DisplayAltAsFigcaption() {
	htmlFormatter.doc.Find("img[alt!='']").Each(func(i int, image *goquery.Selection) {
		alt, _ := image.Attr("alt")

		figure := html.Node{
			Type: html.ElementNode,
			Data: "figure",
		}

		image.WrapNode(&figure)

		figcaption := html.Node{
			Type: html.ElementNode,
			Data: "figcaption",
		}
		figcaption.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: alt,
		})

		image.AfterNodes(&figcaption)
	})
}
