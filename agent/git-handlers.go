package main

import (
	"fmt"

	"net/http"

	"github.com/go-playground/webhooks/v6/bitbucket"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/go-playground/webhooks/v6/gitlab"
)

func (app *application) githubHandler(w http.ResponseWriter, r *http.Request) {
	var release github.PushPayload
	var url, repoName string
	hook, _ := github.New()
	payload, err := hook.Parse(r, github.PushEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}

	switch payload.(type) {

	case github.PushPayload:
		release = payload.(github.PushPayload)
		// Do whatever you want from here...
		fmt.Printf("%s\n", release.Repository.Name)
		url = release.Repository.HTMLURL
		repoName = release.Repository.Name
		app.models.DB.GitPublish(url, repoName, w, r)

	}

}

func (app *application) gitlabHandler(w http.ResponseWriter, r *http.Request) {

	var release gitlab.PushEventPayload
	var url, repoName string
	hook, _ := gitlab.New()
	payload, err := hook.Parse(r, gitlab.PushEvents)
	if err != nil {
		if err == gitlab.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}

	switch payload.(type) {

	case gitlab.PushEventPayload:
		release = payload.(gitlab.PushEventPayload)
		// Do whatever you want from here...
		fmt.Printf("%s\n", release.Repository.Name)
		url = release.Repository.GitHTTPURL
		repoName = release.Repository.Name
		app.models.DB.GitPublish(url, repoName, w, r)

	}
}

func (app *application) bitBucketHandler(w http.ResponseWriter, r *http.Request) {

	var release bitbucket.RepoPushPayload
	var url, repoName string
	hook, _ := bitbucket.New()
	payload, err := hook.Parse(r, bitbucket.RepoPushEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}

	switch payload.(type) {

	case bitbucket.RepoPushPayload:
		release = payload.(bitbucket.RepoPushPayload)
		// Do whatever you want from here...
		fmt.Printf("url url %s\n", release.Repository.Website)
		url = release.Repository.Website
		repoName = release.Repository.Name
		app.models.DB.GitPublish(url, repoName, w, r)

	}
}
