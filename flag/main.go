package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	var build = "0"
	var Version string

	if Version == "" {
		Version = fmt.Sprintf("0.1.%s", build)
	}

	fmt.Println("Version: ", Version)

	plugin := &Plugin{
		Config: &Config{},
		Writer: os.Stdout,
	}

	flag.StringVar(&plugin.Config.branch, "branch", "", "The branch name, ex: staging or master.")
	flag.StringVar(&plugin.Config.inputFile, "input-file", "", "The config file put on drone server.")
	flag.StringVar(&plugin.Config.outputFile, "output-file", "", "Convert secret on the drone server to config file.")
	flag.StringVar(&plugin.Config.repo, "repo", "", "The Repo name on drone Ex:honestbee/<name>")
	flag.BoolVar(&plugin.Config.update, "update", false, "Update the current config on drone server")

	flag.Parse()

	plugin.Exec()
}
