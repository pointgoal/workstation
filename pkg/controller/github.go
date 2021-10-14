package controller

import (
	"context"
	"github.com/google/go-github/v39/github"
	"github.com/pointgoal/workstation/pkg/repository"
	"golang.org/x/oauth2"
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

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken.Token,
			TokenType:   "token",
		},
	)
	client := github.NewClient(oauth2.NewClient(context.Background(), ts))
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
