package gh

import (
	"context"
	"fmt"

	"PromptDefender-Keep/logger"
	"PromptDefender-Keep/score"

	"github.com/google/go-github/v66/github"
	"go.uber.org/zap"
)

type FileHandler struct {
	scorer score.Scorer
	client *github.Client
}

func NewFileHandler(scorer score.Scorer, client *github.Client) *FileHandler {
	return &FileHandler{
		scorer: scorer,
		client: client,
	}
}

func (fh *FileHandler) ShouldRun(ctx context.Context, owner, repo string, prNumber int, promptFiles []string) (bool, error) {
	return true, nil
}

func (fh *FileHandler) RunFilesThroughScoreEndpoint(ctx context.Context, owner, repo, branch string, prNumber int, promptFiles []string) ([]score.PromptScore, error) {
	files, _, err := fh.client.PullRequests.ListFiles(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("error listing files: %w", err)
	}

	var results []score.PromptScore

	for _, file := range files {
		isFileInPromptFiles := false

		for _, promptFile := range promptFiles {
			if file.GetFilename() == promptFile {
				isFileInPromptFiles = true
				break
			}
		}

		if isFileInPromptFiles == false {
			continue
		}

		logger.GetLogger().Info("Processing file", zap.String("filename", file.GetFilename()), zap.Int("pr_number", prNumber), zap.String("branch", branch))

		content, _, _, err := fh.client.Repositories.GetContents(ctx, owner, repo, file.GetFilename(), &github.RepositoryContentGetOptions{
			Ref: branch,
		})

		if err != nil {
			return nil, fmt.Errorf("error getting file content: %w", err)
		}

		if content != nil {
			prompt, err := content.GetContent()
			if err != nil {
				return nil, fmt.Errorf("error getting file content: %w", err)
			}

			logger.GetLogger().Info("Scoring prompt", zap.String("prompt", prompt))

			scoreResult, err := fh.scorer.Score(prompt)
			if err != nil {
				return nil, fmt.Errorf("error scoring prompt: %w", err)
			}

			results = append(results, *scoreResult)
		} else {
			return nil, fmt.Errorf("error getting file content: content is nil")
		}
	}

	return results, nil
}
