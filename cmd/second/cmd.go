package second

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var (
	Cmd = &cli.Command{
		Name:   "second",
		Usage:  "dummy command for example",
		Action: DummyAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from file path`",
				DefaultText: "./cmd/second/config/local.yaml",
				Value:       "./cmd/second/config/local.yaml",
				Required:    false,
			},
		},
	}
)

func DummyAction(cmdCLI *cli.Context) error {
	cfgPath := cmdCLI.String("config")
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to load config from path\"%s\": %w", cfgPath, err)
	}
	fmt.Println("cfg:", cfg)
	fmt.Println("just dummy command for example")
	return nil
}
