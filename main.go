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
var clonesdir string

func init() {
	log.SetLevelFromEnv("WEBHOOK_LOG_LEVEL")

	logger = log.New("webhooks/gitlab/main")

	clonesdir = os.Getenv("WEBHOOK_CLONES_DIR")
	if clonesdir == "" {
		var err error
		clonesdir, err = os.Getwd()
		if err != nil {
			logger.CloneWith(map[string]interface{}{
				"error": err,
			}).Fatal("error getting current working directory for WEBHOOK_CLONES_DIR")
		}
	}
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

	logger.Debug("unmarshaling request body")

	var ev PushEvent
	err = json.Unmarshal(body, &ev)
	if err != nil {
		logger.CloneWith(map[string]interface{}{
			"error": err,
		}).Error("error unmarshaling request body")

		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	logger = logger.CloneWith(map[string]interface{}{
		"event": ev,
	})
	logger.Debug("got event")

	clonedir := fmt.Sprintf("%v/%v", clonesdir, uuid.New())

	logger = logger.CloneWith(map[string]interface{}{
		"clonedir": clonedir,
	})

	logger.Debug("cloning repo")

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

	// generate pipelines

	// send event on queue
}
