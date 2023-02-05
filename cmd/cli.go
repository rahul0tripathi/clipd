package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "initialize the clipd setup",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "keydir",
						Aliases: []string{"dir"},
						Usage:   "define the directory for storing key permutations",
					},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("service initaluzed ")
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
					fmt.Println("server init called")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
