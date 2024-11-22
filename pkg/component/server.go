package component

import (
	"context"

	platformConf "github.com/kiga-hub/arc/conf"
	"github.com/kiga-hub/arc/logging"
	"github.com/kiga-hub/arc/micro"
	"github.com/pangpanglabs/echoswagger/v2"

	"github.com/kiga-hub/data-transmission/pkg/api"
	"github.com/kiga-hub/data-transmission/pkg/upgrade"
)

// DataTransmissionElementKey is Element Key for devops data transmission
var DataTransmissionElementKey = micro.ElementKey("github.com/kiga-hub/dataTransmission")

// DataTransmissionComponent is Component for devops data transmission
type DataTransmissionComponent struct {
	micro.EmptyComponent
	logger  logging.ILogger
	api     *api.Server
	upgrade *upgrade.Client
}

// Name of the component
func (c *DataTransmissionComponent) Name() string {
	return "DataTransmissionComponent"
}

// PreInit called before Init()
func (c *DataTransmissionComponent) PreInit(ctx context.Context) error {
	upgrade.SetDefaultConfig()
	// load config
	return nil
}

// Init the component
func (c *DataTransmissionComponent) Init(server *micro.Server) (err error) {
	c.logger = server.GetElement(&micro.LoggingElementKey).(logging.ILogger)

	if c.upgrade, err = upgrade.New(
		upgrade.WithLogger(c.logger),
		upgrade.WithConfig(upgrade.GetConfig()),
	); err != nil {
		return err
	}

	c.api = api.New(
		api.WithLogger(c.logger),
		api.WithUpgrade(c.upgrade),
	)
	return err
}

// SetDynamicConfig called when dynamic config changed
func (c *DataTransmissionComponent) SetDynamicConfig(nf *platformConf.NodeConfig) error {
	if nf == nil {
		return nil
	}

	return nil
}

// OnConfigChanged called when dynamic config changed
func (c *DataTransmissionComponent) OnConfigChanged(*platformConf.NodeConfig) error {
	return micro.ErrNeedRestart
}

// SetupHandler of echo if the component need
func (c *DataTransmissionComponent) SetupHandler(root echoswagger.ApiRoot, base string) error {
	root.Echo().Static("/", "./ui/dist")
	root.Echo().Static("/static", "./ui/dist/static")
	c.api.Setup(root, base)
	return nil
}

// Start the component
func (c *DataTransmissionComponent) Start(ctx context.Context) error {
	go c.api.Start(ctx)

	return nil
}

// Stop the component
func (c *DataTransmissionComponent) Stop(ctx context.Context) error {
	return nil
}
