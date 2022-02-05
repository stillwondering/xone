package main

import (
	"context"
	"errors"

	"github.com/stillwondering/xone"
	"github.com/stillwondering/xone/sqlite"
	"github.com/urfave/cli/v2"
)

func userCommand(ctx context.Context) *cli.Command {
	cmd := cli.Command{
		Name:    "user",
		Aliases: []string{"u"},
		Usage:   "modify users",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "db",
				Usage:       "use database `FILE`",
				Value:       "xone.db",
				DefaultText: "xone.db",
			},
		},
		Subcommands: []*cli.Command{
			createUserCommand(ctx),
		},
		Before: func(c *cli.Context) error {
			dbfile := c.String("db")
			db := sqlite.NewDB(dbfile)
			if err := db.Open(); err != nil {
				return err
			}

			c.Context = context.WithValue(c.Context, "db", db)
			return nil
		},
		After: func(c *cli.Context) error {
			db := c.Context.Value("db").(*sqlite.DB)
			return db.Close()
		},
	}

	return &cmd
}

func createUserCommand(ctx context.Context) *cli.Command {
	cmd := cli.Command{
		Name:      "create",
		Aliases:   []string{"c"},
		Usage:     "create a new user",
		ArgsUsage: "[email] [password]",
		Action: func(c *cli.Context) error {
			username := c.Args().Get(0)
			if username == "" {
				return errors.New("username missing")
			}
			password := c.Args().Get(1)
			if password == "" {
				return errors.New("password missing")
			}

			db := c.Context.Value("db").(*sqlite.DB)

			userService, err := sqlite.NewUserService(db)
			if err != nil {
				return err
			}

			_, err = userService.Create(ctx, xone.CreateUserData{
				Email:    username,
				Password: password,
			})

			return err
		},
	}

	return &cmd
}
