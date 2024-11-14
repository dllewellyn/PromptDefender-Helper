package gh

import (
	"context"
	"fmt"

	"PromptDefender-Keep/score"

	"github.com/google/go-github/v66/github"
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

func (fh *FileHandler) RunFilesThroughScoreEndpoint(ctx context.Context, owner, repo string, prNumber int) ([]score.PromptScore, error) {
	files, _, err := fh.client.PullRequests.ListFiles(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("error listing files: %w", err)
	}

	var results []score.PromptScore

	for _, file := range files {
		content, _, _, err := fh.client.Repositories.GetContents(ctx, owner, repo, file.GetFilename(), &github.RepositoryContentGetOptions{
			Ref: file.GetSHA(),
		})
		if err != nil {
			return nil, fmt.Errorf("error getting file content: %w", err)
		}

		if content != nil {
			prompt, err := content.GetContent()
			if err != nil {
				return nil, fmt.Errorf("error getting file content: %w", err)
			}

			scoreResult, err := fh.scorer.Score(prompt)
			if err != nil {
				return nil, fmt.Errorf("error scoring prompt: %w", err)
			}

			results = append(results, *scoreResult)
		}
	}

	return results, nil
}
