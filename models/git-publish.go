package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	billy "github.com/go-git/go-billy/v5"
	memfs "github.com/go-git/go-billy/v5/memfs"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	htt "github.com/go-git/go-git/v5/plumbing/transport/http"
	memory "github.com/go-git/go-git/v5/storage/memory"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const (
	streamName     = "GITMETRICS"
	streamSubjects = "GITMETRICS.*"
	eventSubject   = "GITMETRICS.event"
	allSubject     = "GITMETRICS.all"
	version        = "1.0.0"
)

var storer *memory.Storage
var fs billy.Filesystem

type DBModel struct {
	JS       nats.JetStreamContext
	gituser  string
	gittoken string
}
type Branches []string

func (m *DBModel) GitPublish(url, repoName string, w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New()
	metrics := Gitevent{
		Repository: repoName,
		Uuid:       uuid,
	}
	storer = memory.NewStorage()
	fs = memfs.New()
	// Authentication
	auth := &htt.BasicAuth{
		Username: m.gituser,
		Password: m.gittoken,
	}
	repClone, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:  url,
		Auth: auth,
	})
	if err != nil {
		log.Fatal(err)
	}
	remote, err := repClone.Remote("origin")
	if err != nil {
		panic(err)
	}
	refList, err := remote.List(&git.ListOptions{})
	if err != nil {
		panic(err)
	}

	refPrefix := "refs/heads/"

	var branches Branches
	for _, ref := range refList {

		refName := ref.Name().String()
		if !strings.HasPrefix(refName, refPrefix) {
			continue
		}
		branchName := refName[len(refPrefix):]

		branches = append(branches, branchName)

	}
	metrics.Availablebranches = totalbranches(&branches)
	// ... retrieving the branch being pointed by HEAD
	ref, err := repClone.Head()
	if err != nil {
		panic(err)
	}
	// ... retrieving the commit object
	commit, err := repClone.CommitObject(ref.Hash())
	if err != nil {
		panic(err)
	}

	metrics.CommitedBy = commit.Author.Name
	metrics.CommitedAt = commit.Author.When
	metrics.Commitmessage = commit.Message

	stats, _ := commit.Stats()

	metrics.Commitstat = getCommitStats(stats)

	metricsJson, _ := json.Marshal(metrics)
	_, err = m.JS.Publish(eventSubject, metricsJson)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(metricsJson))
	log.Printf("Metrics with url:%s has been published\n", url)

}
func totalbranches(b *Branches) string {
	var sb strings.Builder
	for _, bran := range *b {
		sb.WriteString(bran)
		sb.WriteString(",")
	}
	return sb.String()
}
func getCommitStats(stat object.FileStats) string {
	var sb strings.Builder
	for _, comm := range stat {
		sb.WriteString(comm.Name)
		sb.WriteString(",")
		sb.WriteString("Add" + ":")
		sb.WriteString(fmt.Sprintf("%v", comm.Addition))
		sb.WriteString(",")
		sb.WriteString("Del" + ":")
		sb.WriteString(fmt.Sprintf("%v", comm.Deletion))
		sb.WriteString("  ")
	}
	return sb.String()
}
