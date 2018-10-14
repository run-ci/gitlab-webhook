package main

// PushEvent is a Gitlab push event.
type PushEvent struct {
	Repository Repository `json:"repository"`
}

// Repository is a Gitlab repository.
type Repository struct {
	GitHTTPURL string `json:"git_http_url"`
}
