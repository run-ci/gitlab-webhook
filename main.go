package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/uuid"

	"gitlab.com/run-ci/webhooks/gitlab/log"

	git "gopkg.in/src-d/go-git.v4"
)

var logger *log.Logger

func init() {
	log.SetLevelFromEnv("WEBHOOK_LOG_LEVEL")

	logger = log.New("webhooks/gitlab/main")
}

func main() {
	logger.Info("booting gitlab webhook")

	http.HandleFunc("/", handle)

	http.ListenAndServe(":9090", nil)
}

func handle(wr http.ResponseWriter, req *http.Request) {
	logger = logger.CloneWith(map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
	})

	logger.Debug("handling request")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.CloneWith(map[string]interface{}{
			"error": err,
		}).Error("error reading request body")

		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ev PushEvent
	err = json.Unmarshal(body, &ev)
	if err != nil {
		logger.CloneWith(map[string]interface{}{
			"error": err,
		}).Error("error unmarshaling request body")

		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debugf("got event: %+v", ev)

	basedir, err := os.Getwd()
	if err != nil {
		logger.CloneWith(map[string]interface{}{
			"error": err,
			"event": ev,
		}).Error("error getting working directory")
	}

	clonedir := fmt.Sprintf("%v/%v", basedir, uuid.New())

	_, err = git.PlainClone(clonedir, false, &git.CloneOptions{
		URL:      ev.Repository.GitHTTPURL,
		Progress: os.Stdout,
	})
	if err != nil {
		logger.CloneWith(map[string]interface{}{
			"error": err,
			"event": ev,
		}).Error("unable to clone git repository")
	}

	// read pipeline files

	// generate pipelines

	// send event on queue
}
