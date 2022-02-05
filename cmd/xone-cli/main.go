package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx, os.Args); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	appName := filepath.Base(args[0])

	app := cli.App{
		Name:  appName,
		Usage: "xone's command line management tool",
		Commands: []*cli.Command{
			userCommand(ctx),
		},
	}

	return app.Run(args)
}
