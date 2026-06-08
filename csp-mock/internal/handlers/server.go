package handlers

import (
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/util"
)

// Server holds all application dependencies and owns route registration.
// Add new repositories to Repos as the API grows.
type Server struct {
	repos *repository.Repos
	clock util.Clock
}

func NewServer(repos *repository.Repos, clock util.Clock) *Server {
	return &Server{repos: repos, clock: clock}
}
