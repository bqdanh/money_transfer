package cmd

import (
	"github.com/bqdanh/money_transfer/cmd/second"
	"github.com/bqdanh/money_transfer/cmd/start_server"
	"github.com/urfave/cli/v2"
)

func AppCommandLineInterface() *cli.App {
	appCli := cli.NewApp()
	appCli.Action = start_server.StartServerAction
	appCli.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "Load configuration from file path`",
			DefaultText: "./configs/server/local.yaml",
			Value:       "./configs/server/local.yaml",
			Required:    false,
		},
	}

	appCli.Commands = []*cli.Command{
		start_server.Cmd,
		second.Cmd,
	}

	return appCli
}
