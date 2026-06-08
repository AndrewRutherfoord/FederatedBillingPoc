package handlers

import (
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/scheduler"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/util"
)

// Server holds all application dependencies and owns route registration.
// Add new repositories to Repos as the API grows.
type Server struct {
	repos      *repository.Repos
	clock      util.Clock
	scheduler  *scheduler.Scheduler
}

func NewServer(repos *repository.Repos, clock util.Clock, sched *scheduler.Scheduler) *Server {
	return &Server{repos: repos, clock: clock, scheduler: sched}
}
