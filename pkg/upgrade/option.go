package upgrade

import (
	"github.com/kiga-hub/arc/logging"
)

// Option is a function that will set up option.
type Option func(opts *Client)

func loadOptions(options ...Option) *Client {
	opts := &Client{}
	for _, option := range options {
		option(opts)
	}
	if opts.logger == nil {
		opts.logger = new(logging.NoopLogger)
	}
	if opts.config == nil {
		opts.config = GetConfig()
	}
	return opts
}

// WithLogger -
func WithLogger(logger logging.ILogger) Option {
	return func(opts *Client) {
		opts.logger = logger
	}
}

// WithConfig -
func WithConfig(c *Config) Option {
	return func(opts *Client) {
		opts.config = c
	}
}
