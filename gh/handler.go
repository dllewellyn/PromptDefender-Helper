package gh

import (
	"PromptDefender-Keep/logger"
	"PromptDefender-Keep/score"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	ghinstallation "github.com/bradleyfalzon/ghinstallation/v2"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v66/github"
)

func HandleWebhook(c *gin.Context, scorer score.Scorer) {
	logger.GetLogger().Info("Received webhook", zap.Any("request", c.Request.Header), zap.Any("body", c.Request.Body))

	payload, err := github.ValidatePayload(c.Request, []byte(os.Getenv("GITHUB_WEBHOOK_SECRET")))
	if err != nil {
		logger.GetLogger().Error("Could not validate payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if c.GetHeader("X-GitHub-Event") != "pull_request" {
		logger.GetLogger().Error("Invalid event type", zap.String("event", c.GetHeader("X-GitHub-Event")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event type"})
	}

	event := github.PullRequestEvent{}
	if err := json.Unmarshal(payload, &event); err != nil {
		logger.GetLogger().Error("Could not parse webhook", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse webhook"})
		return
	}

	if event.GetAction() == "opened" || event.GetAction() == "synchronize" || event.GetAction() == "reopened" {
		go handlePullRequest(event, scorer)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func handlePullRequest(event github.PullRequestEvent, scorer score.Scorer) ([]score.PromptScore, error) {
	ctx := context.Background()

	owner := event.Repo.Owner.GetLogin()
	repo := event.Repo.GetName()
	prNumber := event.GetNumber()
	branchName := event.PullRequest.Head.GetRef()

	// Get the latest version of the .github/prompt-defender.yml file from the branch
	branch := event.PullRequest.Head.GetRef()

	installationID := event.Installation.GetID()
	githubAppId := os.Getenv(("GITHUB_APPLICATION_ID"))

	githubAppIdInt, err := strconv.Atoi(githubAppId)
	if err != nil {
		fmt.Printf("Error converting GITHUB_APPLICATION_ID to int: %v\n", err)
		return nil, err
	}

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, int64(githubAppIdInt), installationID, "github-app.pem")

	if err != nil {
		fmt.Printf("Error creating github client: %v\n", err)
	}

	client := github.NewClient(&http.Client{Transport: itr})

	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, ".github/prompt-defender.yml", &github.RepositoryContentGetOptions{Ref: branch})
	if err != nil {
		fmt.Printf("Error getting file content from GitHub: %v\n", err)
		return nil, err
	}

	content, err := fileContent.GetContent()
	if err != nil {
		fmt.Printf("Error decoding file content: %v\n", err)
		return nil, err
	}

	config, err := LoadConfigFromString(content)

	if err != nil {
		fmt.Printf("Error loading config from string: %v\n", err)
		return nil, err
	}

	fileHandler := NewFileHandler(scorer, client)

	logger.GetLogger().Info("Checking if should run", zap.Any("config", config), zap.Any("owner", owner), zap.Any("repo", repo), zap.Any("prNumber", prNumber))

	if shouldRun, err := fileHandler.ShouldRun(ctx, owner, repo, prNumber, config.Prompts); err == nil && shouldRun {

		response, err := fileHandler.RunFilesThroughScoreEndpoint(ctx, owner, repo, branchName, prNumber, config.Prompts)

		if err != nil {
			fmt.Printf("Error running files through score endpoint: %v\n", err)
		}

		return response, nil
	}

	return make([]score.PromptScore, 0), nil

}

func init() {
	requiredEnvVars := []string{"GITHUB_WEBHOOK_SECRET", "GITHUB_APPLICATION_ID"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			logger.GetLogger().Fatal(fmt.Sprintf("%s environment variable not set", envVar))
		}
	}
}
