package main

import (
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/urfave/cli/v2"
	"golang.design/x/hotkey/mainthread"
	"log"
	"os"
)

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "initialise the clipd setup",
				Action: func(cCtx *cli.Context) error {
					initd()
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start the clipd server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "keydir",
						Aliases: []string{"dir"},
						Usage:   "define the directory where key permutations are stored",
					},

					&cli.StringFlag{
						Name:  "salt",
						Usage: "define the directory for storing key permutations",
					},
				},
				Action: func(cCtx *cli.Context) error {
					salt := cCtx.String("salt")
					if salt == "" {
						return errors.New("salt not given")
					}
					keyDir := cCtx.String("keydir")
					if keyDir == "" {
						return errors.New("key directory not given")
					}
					mainthread.Init(func() {
						start(salt, keyDir)
					})
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
