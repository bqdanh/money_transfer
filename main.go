package main

import (
	"os"

	"github.com/bqdanh/money_transfer/cmd"
)

func main() {
	appCli := cmd.AppCommandLineInterface()
	if err := appCli.Run(os.Args); err != nil {
		panic(err)
	}
}
