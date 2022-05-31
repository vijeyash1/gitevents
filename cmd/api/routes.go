package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/webhooks", app.githubHandler)
	router.HandlerFunc(http.MethodGet, "/v1/gitlab", app.gitlabHandler)
	router.HandlerFunc(http.MethodGet, "/v1/bitbucket", app.bitBucketHandler)

	return router
}
