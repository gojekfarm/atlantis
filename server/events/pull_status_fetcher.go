package events

import "github.com/gojekfarm/atlantis/server/events/models"

// PullStatusFetcher fetches our internal model of a pull requests status
type PullStatusFetcher interface {
	GetPullStatus(pull models.PullRequest) (*models.PullStatus, error)
}
