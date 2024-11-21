package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/kiga-hub/arc/micro"
	basicComponent "github.com/kiga-hub/arc/micro/component"
	"github.com/spf13/cobra"

	"github.com/kiga-hub/data-transmission/pkg/component"
)

// ServerCmd -
var ServerCmd = &cobra.Command{
	Use:   "run",
	Short: "run dev opts data transmission server",
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(ServerCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	// recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", "recover", r)
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	components := []micro.IComponent{
		&basicComponent.LoggingComponent{},
		// &mysql.Component{},
	}
	components = append(components, &component.DataTransmissionComponent{})

	server, err := micro.NewServer(
		AppName,
		AppVersion,
		components,
	)
	if err != nil {
		panic(err)
	}
	err = server.Init()
	if err != nil {
		panic(err)
	}

	err = server.Run()
	if err != nil {
		panic(err)
	}
}
