package main

import (
	"fmt"

	"github.com/go-playground/webhooks/v6/github"

	"net/http"
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

	fmt.Fprintf(w, " this is gitlab handler")
}

func (app *application) bitBucketHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, " this is bitbucket handler")
}
