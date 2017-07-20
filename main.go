package main

import "github.com/mgutz/ansi"
import "github.com/briandowns/spinner"
import "fmt"
import "os"
import "bufio"
import "io/ioutil"
import "time"
import "github.com/Medium/medium-sdk-go"
import "github.com/google/go-github/github"
import "golang.org/x/oauth2"
import "context"

func main() {
	logger := NewPrettyLogger()

	config, err := LoadConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load config file: %s", err))
	}

	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	var configChanged bool

	if config.MediumAccessToken == "" {
		mediumAccessToken := promptForMediumAccessToken()
		config.MediumAccessToken = mediumAccessToken
		configChanged = true
	}

	mediumClient := medium.NewClientWithAccessToken(config.MediumAccessToken)

	if config.MediumUserId == "" {
		logger.Info("Fetching Medium user")

		s.Start()
		user, err := mediumClient.GetUser("")
		s.Stop()

		if err != nil {
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to load Medium user: %s", err))
			}
		}

		config.MediumUserId = user.ID
		configChanged = true
	}

	if config.GitHubAccessToken == "" {
		if githubAccessToken := promptForGitHubAccessToken(); githubAccessToken != "" {
			config.GitHubAccessToken = githubAccessToken
			configChanged = true
		}
	}

	if configChanged {
		if err := SaveConfig(config); err != nil {
			logger.Error(fmt.Sprintf("Failed to save config file: %s", err))
		}
	}

	markdownPath := os.Args[1]

	if _, err := os.Stat(markdownPath); os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("File does not exist: %s", markdownPath))
	}

	markdownContent, err := ioutil.ReadFile(markdownPath)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read %s: %s", markdownPath, err))
	}

	ctx := context.Background()
	githubClient := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubAccessToken},
	)))

	metadata, markdownContent, err := ExtractMarkdownMetadata(markdownContent)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to extract metadata from Markdown file: %s", err))
	}

	postHtml := string(MarkdownToHTML(markdownContent)[:])

	htmlFormatter, _ := NewHtmlFormatter(markdownPath, postHtml, logger, ctx)

	if err := htmlFormatter.ReplaceCodeBlocks(githubClient); err != nil {
		logger.Error(fmt.Sprintf("Error replacing code blocks with gists: %s", err))
	}

	if err := htmlFormatter.ReplaceImages(mediumClient); err != nil {
		logger.Error(fmt.Sprintf("Error replacing images: %s", err))
	}

	formattedHtml, err := htmlFormatter.Html()
	if err != nil {
		logger.Error(fmt.Sprintf("Error formatting HTML: %s", err))
	}

	logger.Info("Creating Medium post")
	s.Start()
	post, err := mediumClient.CreatePost(medium.CreatePostOptions{
		UserID:        config.MediumUserId,
		Content:       formattedHtml,
		ContentFormat: medium.ContentFormatHTML,
		PublishStatus: medium.PublishStatusDraft,
		Tags:          metadata.Tags,
	})
	s.Stop()

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create Medium post: %s", err))
	}

	fmt.Println("")
	logger.Success(fmt.Sprintf("Successfully created Medium post.\n  You can find it here: %s%s%s", ansi.ColorCode("white+u"), post.URL, ansi.ColorCode("reset")))
}

func printUsage() {
	fmt.Println("")
	fmt.Println("Usage: md2markdown MARKDOWN_FILE_PATH")
	fmt.Println("")
}

func promptForGitHubAccessToken() string {
	fmt.Println("")
	fmt.Println("If you wish to syntax highlight code snippets, you will need to provide a GitHub access token.")
	fmt.Println("This token will only be used to create private gists which will be embedded within your post.")
	fmt.Println("")
	fmt.Printf("You can generate a personal access token here: %shttps://github.com/settings/tokens%s\n", ansi.ColorCode("white+u"), ansi.ColorCode("reset"))
	fmt.Println("")
	fmt.Println("Enter personal access token or hit enter to skip:")
	fmt.Println("")
	fmt.Printf("%s>%s ", ansi.ColorCode("white+b"), ansi.ColorCode("reset"))

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	githubAccessToken := scanner.Text()

	fmt.Println("")

	return githubAccessToken
}

func promptForMediumAccessToken() string {
	fmt.Println("")
	fmt.Println("In order to create the post on Medium, you will need to provide a Medium access token.")
	fmt.Println("")
	fmt.Printf("You can generate an integration token here: %shttps://medium.com/me/settings%s\n", ansi.ColorCode("white+u"), ansi.ColorCode("reset"))
	fmt.Println("")
	fmt.Println("Enter integration token:")
	fmt.Println("")
	fmt.Printf("%s>%s ", ansi.ColorCode("white+b"), ansi.ColorCode("reset"))

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	mediumAccessToken := scanner.Text()

	if mediumAccessToken == "" {
		return promptForMediumAccessToken()
	}

	fmt.Println("")

	return mediumAccessToken
}
