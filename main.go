package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"

	"gitlab.com/run-ci/webhooks/gitlab/log"
	"gitlab.com/run-ci/webhooks/gitlab/queue"
	pkg "gitlab.com/run-ci/webhooks/pkg.git"

	git "gopkg.in/src-d/go-git.v4"
	yaml "gopkg.in/yaml.v2"
)

var logger *log.Logger
var clonesdir string
var q pkg.PipelineSender

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

	q = queue.NewEchoQueue()
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
		}).Error("unable to clone git repository")

		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	pipelinesdir := fmt.Sprintf("%v/pipelines", clonedir)
	finfos, err := ioutil.ReadDir(pipelinesdir)
	if err != nil {
		logger.CloneWith(map[string]interface{}{
			"error": err,
			"dir":   pipelinesdir,
		}).Error("unable to list directory")

		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, finfo := range finfos {
		// The file name corresponds to the pipeline name.
		name := strings.Split(finfo.Name(), ".")[0]

		logger := logger.CloneWith(map[string]interface{}{
			"pipeline_name": name,
		})

		logger.Debug("processing pipeline")

		logger.Debug("opening pipeline file")

		f, err := os.Open(fmt.Sprintf("%v/%v", pipelinesdir, finfo.Name()))
		if err != nil {
			logger.CloneWith(map[string]interface{}{
				"error": err,
			}).Error("unable to open pipeline file")

			continue
		}

		logger.Debug("reading pipeline file")

		buf, err := ioutil.ReadAll(f)
		if err != nil {
			logger.CloneWith(map[string]interface{}{
				"error": err,
			}).Error("unable to read pipeline file")

			continue
		}

		logger.Debug("loading pipeline")

		var p pkg.Pipeline
		err = yaml.UnmarshalStrict(buf, &p)
		if err != nil {
			logger = logger.CloneWith(map[string]interface{}{
				"error": err,
			})
			logger.Error("unable to unmarshal pipeline")

			continue
		}

		p.Name = name
		p.Remote = ev.Repository.GitHTTPURL

		logger = logger.CloneWith(map[string]interface{}{
			"pipeline": p,
		})
		logger.Debug("sending pipeline")

		err = q.SendPipeline(p)
		if err != nil {
			logger.CloneWith(map[string]interface{}{
				"error": err,
			}).Error("unable to send pipeline")

			continue
		}
	}

	wr.WriteHeader(http.StatusNoContent)
	return
}
