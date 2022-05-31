package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/github", app.githubHandler)
	router.HandlerFunc(http.MethodPost, "/v1/gitlab", app.gitlabHandler)
	router.HandlerFunc(http.MethodPost, "/v1/bitbucket", app.bitBucketHandler)

	return router
}
