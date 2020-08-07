package data

import (
	"context"
	"github.com/google/go-github/github"
	"time"
)

type RibsRepo struct {
	Id              int64
	HtmlUrl         string
	FullName        string
	Description     string
	Homepage        string
	ForksCount      int
	StargazersCount int
	Languages       map[string]int
}

type RibsIssue struct {
	HtmlUrl  string
	Number   int
	Title    string
	State    string
	Labels   []RibsIssueLabel
	Comments []RibsIssueComment
}

type RibsIssueComment struct {
	Id            int64
	HtmlUrl       string
	UserName      string
	UserAvatarUrl string
	Created       time.Time
	Updated       time.Time
	Body          string
}

type RibsIssueLabel struct {
	Id    int64
	Url   string
	Name  string
	Color string
}

func FetchRepo(owner string, repo string) (RibsRepo, error) {
	client := github.NewClient(nil)
	var parsedRepo RibsRepo
	ghRepo, _, ghRepoErr := client.Repositories.Get(context.Background(), owner, repo)
	if ghRepoErr != nil {
		return parsedRepo, ghRepoErr
	} else {
		parsedRepo = parseRepo(*ghRepo)
	}
	ghLangs, _, ghLangsErr := client.Repositories.ListLanguages(context.Background(), owner, repo)
	if ghLangsErr != nil {
		return parsedRepo, ghLangsErr
	} else {
		parsedRepo.Languages = ghLangs
	}
	return parsedRepo, nil
}

func FetchOpenIssues(owner string, repo string) ([]RibsIssue, error) {
	client := github.NewClient(nil)
	issuesOpts := github.IssueListByRepoOptions{
		State:       "open",
		Sort:        "comments",
		Direction:   "desc",
		ListOptions: github.ListOptions{},
	}
	var parsedIssues []RibsIssue
	ghIssues, _, ghIssuesErr := client.Issues.ListByRepo(context.Background(), owner, repo, &issuesOpts)
	if ghIssuesErr != nil {
		return nil, ghIssuesErr
	} else {
		parsedIssues = parseIssues(ghIssues)
	}
	var commentOpts = github.IssueListCommentsOptions{}
	var comments []RibsIssueComment
	for _, e := range parsedIssues {
		cs, _, _ := client.Issues.ListComments(context.Background(), owner, repo, e.Number, &commentOpts)
		comments = append(comments, parseComments(cs)...)
	}
	return parsedIssues, nil
}

func parseRepo(repo github.Repository) RibsRepo {
	return RibsRepo{
		Id:              *repo.ID,
		HtmlUrl:         *repo.HTMLURL,
		FullName:        *repo.FullName,
		Description:     *repo.Description,
		Homepage:        *repo.Homepage,
		ForksCount:      *repo.ForksCount,
		StargazersCount: *repo.StargazersCount,
		Languages:       nil,
	}
}

func parseIssues(issues []*github.Issue) []RibsIssue {
	var parsedIssues []RibsIssue
	for _, e := range issues {
		parsedIssues = append(parsedIssues, RibsIssue{
			HtmlUrl: *e.HTMLURL,
			Number:  *e.Number,
			Title:   *e.Title,
			State:   *e.State,
			Labels:  parseLabels(e.Labels),
		})
	}
	return parsedIssues
}

func parseLabels(labels []github.Label) []RibsIssueLabel {
	var parsedLabels []RibsIssueLabel
	for _, e := range labels {
		parsedLabels = append(parsedLabels, RibsIssueLabel{
			Id:    *e.ID,
			Url:   *e.URL,
			Name:  *e.Name,
			Color: *e.Color,
		})
	}
	return parsedLabels
}

func parseComments(comments []*github.IssueComment) []RibsIssueComment {
	var parsedComments []RibsIssueComment
	for _, e := range comments {
		parsedComments = append(parsedComments, RibsIssueComment{
			Id:            *e.ID,
			HtmlUrl:       *e.HTMLURL,
			UserName:      *e.User.Name,
			UserAvatarUrl: *e.User.AvatarURL,
			Created:       *e.CreatedAt,
			Updated:       *e.UpdatedAt,
			Body:          *e.Body,
		})
	}
	return parsedComments
}
