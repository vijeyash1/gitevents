package models

import (
	"fmt"
	"net/http"

	//billy "github.com/go-git/go-billy/v5"
	//memory "github.com/go-git/go-git/v5/storage/memory"
	"github.com/nats-io/nats.go"
)

// var storer *memory.Storage
// var fs billy.Filesystem

type DBModel struct {
	JS nats.JetStreamContext
}

func (m *DBModel) GithubPublish(url, repoName string, w http.ResponseWriter) {
	//fmt.Fprintf(w, "url: %s reponame: %s", url, repoName)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("url:%s repo:%s", url, repoName)))
	fmt.Println("the response is ", url, repoName)

}
