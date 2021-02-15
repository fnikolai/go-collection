package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/fnikolai/go-collection/terminal"
)

var (
	// Version that is passed on compile time through -ldflags
	Version = "built locally"

	// GitCommit that is passed on compile time through -ldflags
	GitCommit = "none"

	// GitBranch that is passed on compile time through -ldflags
	GitBranch = "none"

	// BuildTime that is passed on compile time through -ldflags
	BuildTime = "none"

	// HumanVersion is a human readable app version
	HumanVersion = fmt.Sprintf("%s - %.7s (%s) %s", Version, GitCommit, GitBranch, BuildTime)
)

func main() {
	// create cancelable context
	topContext, cancel := context.WithCancel(context.Background())

	//
	// Configure CLI appearance
	//
	app := &cli.App{
		Name:                 "test terminnal",
		EnableBashCompletion: true,
		Usage:                "demo usage of terminal logger and signal handler",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "verbose, V",
				Usage: "Print verbose output [debug,info,warning,error,fatal,panic]",
				Value: "info",
			},
		},
		Before: func(c *cli.Context) error {
			if err := terminal.SetLogger(c); err != nil {
				return err
			}

			return nil
		},
		Commands: func(TopContext context.Context) []*cli.Command {
			return []*cli.Command{}
		}(topContext),
	}

	terminal.HandleSignals(cancel, func() error {
		return app.Run(os.Args)
	})
}
