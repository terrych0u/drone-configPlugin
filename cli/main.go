package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func run(c *cli.Context) error {
	plugin := Plugin{
		Config: Config{
			branch:     c.String("branch"),
			inputFile:  c.String("input-file"),
			outputFile: c.String("output-file"),
			repo:       c.String("repo"),
			update:     c.Bool("update"),
		},

		Writer: os.Stdout,
	}

	return plugin.Exec()
}

var build = "0"
var Version string

func main() {
	if Version == "" {
		Version = fmt.Sprintf("0.1.%s", build)
	}

	app := cli.NewApp()
	app.Name = "Drone config"
	app.Usage = "Input config file on drone that can output in CI process"
	app.Copyright = "Copyright (c) 2018 Terry Chou"
	app.Authors = []cli.Author{
		{
			Name:  "Terry Chou",
			Email: "tequilas721@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "b, branch",
			Usage: "The branch name, ex: staging or master.",
		},
		cli.StringFlag{
			Name:  "i, input-file",
			Usage: "The config file put on drone server.",
		},
		cli.StringFlag{
			Name:  "o, output-file",
			Usage: "Convert secret on the drone server to config file.",
		},
		cli.StringFlag{
			Name:  "r, repo",
			Usage: "The Repo name on drone `honestbee/<name>`",
		},
		cli.BoolFlag{
			Name:   "update",
			Usage:  "Update the current config on drone server.",
			Hidden: true,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
