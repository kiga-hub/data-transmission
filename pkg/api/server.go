package api

import (
	"context"

	"github.com/kiga-hub/arc/logging"
	"github.com/kiga-hub/data-transmission/pkg/upgrade"
	"github.com/pangpanglabs/echoswagger/v2"
)

// Server -
type Server struct {
	logger  logging.ILogger
	upgrade *upgrade.Client
}

// New -
func New(opts ...Option) *Server {
	srv := loadOptions(opts...)
	return srv
}

// Setup -
func (s *Server) Setup(root echoswagger.ApiRoot, base string) {
	s.setupTransmission(root, base)
}

// Start -
func (s *Server) Start(ctx context.Context) {
}
