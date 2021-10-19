package controller

import (
	"context"
	"fmt"
	"github.com/google/go-github/v39/github"
	"github.com/pointgoal/workstation/pkg/repository"
	"golang.org/x/oauth2"
	"strings"
)

const (
	PerPageDefault = 10
	PageDefault    = 1
)

type Installation struct {
	RepoSource   string        `yaml:"repoSource" json:"repoSource"`
	Organization string        `yaml:"organization" json:"organization"`
	AvatarUrl    string        `yaml:"avatarUrl" json:"avatarUrl"`
	Repos        []*Repository `yaml:"repos" json:"repos"`
}

type Repository struct {
	FullName string `yaml:"fullName" json:"fullName"`
	Name     string `yaml:"name" json:"name"`
}

// ListUserInstallations returns repositories from code repo created by user.
// The repos would have access permission with Github app named as workstation.
func ListUserInstallationsFromGithub(user string) ([]*Installation, error) {
	res := make([]*Installation, 0)

	// 1: Get access code from repo
	repo := repository.GetRepository()

	accessToken, err := repo.GetAccessToken("github", user)
	if err != nil {
		return res, err
	}

	client := getGithubClient(accessToken.Token)
	listOpts := &github.ListOptions{}

	installsFromGithub, _, err := client.Apps.ListUserInstallations(context.Background(), listOpts)
	if err != nil {
		return res, err
	}

	for i := range installsFromGithub {
		install := &Installation{
			AvatarUrl:    installsFromGithub[i].GetAccount().GetAvatarURL(),
			RepoSource:   "github",
			Organization: installsFromGithub[i].GetAccount().GetLogin(),
			Repos:        make([]*Repository, 0),
		}

		// list repositories from this installation
		userReposFromRepo, _, err := client.Apps.ListUserRepos(context.Background(), installsFromGithub[i].GetID(), listOpts)
		if err != nil {
			return res, err
		}

		for j := range userReposFromRepo.Repositories {
			repo := &Repository{
				FullName: userReposFromRepo.Repositories[j].GetFullName(),
				Name:     userReposFromRepo.Repositories[j].GetName(),
			}
			install.Repos = append(install.Repos, repo)
		}
		res = append(res, install)
	}

	return res, nil
}

// ListCommitsFromGithub returns commits from remote github repository.
// The repos would have access permission with Github app named as workstation.
func ListCommitsFromGithub(src *repository.Source, branch, accessToken string, perPage, page int) ([]*Commit, error) {
	res := make([]*Commit, 0)

	client := getGithubClient(accessToken)

	opts := &github.CommitsListOptions{
		SHA: branch,
		ListOptions: github.ListOptions{
			PerPage: perPage,
			Page:    page,
		},
	}

	// repo was stored with format of owner/repo
	srcTokens := strings.Split(src.Repository, "/")
	if len(srcTokens) < 2 {
		return res, fmt.Errorf("invalid repository in DB repo:%s", src.Repository)
	}

	commits, _, err := client.Repositories.ListCommits(context.Background(), srcTokens[0], srcTokens[1], opts)
	if err != nil {
		return res, err
	}

	for i := range commits {
		res = append(res, &Commit{
			Id:           commits[i].GetSHA(),
			Url:          commits[i].GetHTMLURL(),
			Message:      commits[i].GetCommit().GetMessage(),
			Date:         commits[i].GetCommit().GetCommitter().GetDate(),
			Committer:    commits[i].GetCommit().GetCommitter().GetName(),
			CommitterUrl: commits[i].GetCommitter().GetHTMLURL(),
		})
	}

	return res, nil
}

// normalize page and perPage
func normalizePage(perPage, page int) (int, int) {
	if perPage < 0 {
		perPage = PerPageDefault
	}

	if page < 0 {
		page = PageDefault
	}

	return perPage, page
}

// Get github client with accessToken
func getGithubClient(accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
			TokenType:   "token",
		},
	)
	client := github.NewClient(oauth2.NewClient(context.Background(), ts))

	return client
}
