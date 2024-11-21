package api

import (
	"context"

	"github.com/kiga-hub/arc/logging"
	"github.com/kiga-hub/data-transmission/pkg/upgrade"
	"github.com/pangpanglabs/echoswagger/v2"
)

// Server - api接口管理结构
type Server struct {
	logger  logging.ILogger
	upgrade *upgrade.Client
}

// New - 新建
func New(opts ...Option) *Server {
	srv := loadOptions(opts...)
	return srv
}

// Setup - 安装接口
func (s *Server) Setup(root echoswagger.ApiRoot, base string) {
	s.setupUpgrade(root, base)
}

// Start loop of warning auto resolve
func (s *Server) Start(ctx context.Context) {
}
