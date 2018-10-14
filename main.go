package main

import (
	"io/ioutil"
	"net/http"

	"gitlab.com/run-ci/gitlab-webhook/log"
)

var logger *log.Logger

func init() {
	log.SetLevelFromEnv("WEBHOOK_LOG_LEVEL")

	logger = log.New("gitlab-webhook/main")
}

func main() {
	logger.Info("booting gitlab-webhook")

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

		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debugf("got request body: %s", body)
}
