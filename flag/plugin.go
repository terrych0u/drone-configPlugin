package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	bothInputAndOutputExist       = "Error: input and output only one can exist at the same time."
	missingInputConfigFileAndRepo = "Error: can't find input config file or repo name."
	missingBranch                 = "Error: can't find branch name."
)

type (
	// Config for the plugin.
	Config struct {
		branch     string
		inputFile  string
		outputFile string
		repo       string
		update     bool
	}

	// Plugin structure
	Plugin struct {
		Config *Config
		Writer io.Writer
	}
)

func checkEnv() {
	cmd := exec.Command("drone",
		"-t", os.Getenv("DRONE_TOKEN"),
		"-s", os.Getenv("DRONE_SERVER"),
		"info",
	)

	stdoutStderr, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", stdoutStderr)
	if err != nil {
		fmt.Println("Error: please reference http://docs.drone.io/cli-authentication that setup your drone environmental variables.")
		os.Exit(1)
	}
	return
}

func (p *Plugin) addSecret(data string) (string, error) {

	action := "add"
	nPtr := &action

	if p.Config.update {
		*nPtr = "update"
	}

	cmd := exec.Command("drone", "secret", action,
		fmt.Sprintf("-repository %s", p.Config.repo),
		"-event=pull_request",
		"-event=push",
		"-event=tag",
		"-event=deployment",
		fmt.Sprintf("-name=%s_CONFIG_FILE", strings.ToUpper(p.Config.branch)),
		fmt.Sprintf("-value %s", data),
	)

	// fmt.Println(cmd)
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	return "", nil
}

func (p *Plugin) Encode() (string, error) {
	if len(p.Config.repo) == 0 && len(p.Config.inputFile) == 0 {
		fmt.Println(missingInputConfigFileAndRepo)
		os.Exit(1)
	}

	checkEnv()

	buff, err := ioutil.ReadFile(p.Config.inputFile)
	if err != nil {
		return "", err
	}

	data := base64.StdEncoding.EncodeToString(buff)

	p.addSecret(data)

	return "", nil
}

func (p *Plugin) Decode() (string, error) {
	data := os.Getenv(fmt.Sprintf("%s_CONFIG_FILE", strings.ToUpper(p.Config.branch)))

	content, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("Error:", err)
	}

	ioutil.WriteFile(p.Config.outputFile, content, 0644)

	return "", nil
}

func (p *Plugin) Exec() error {
	if len(p.Config.branch) == 0 {
		fmt.Println(missingBranch)
		os.Exit(1)
	}

	if len(p.Config.inputFile) != 0 && len(p.Config.outputFile) != 0 {
		fmt.Println(bothInputAndOutputExist)
		os.Exit(1)
	}

	if len(p.Config.outputFile) != 0 {
		p.Decode()
	} else {
		p.Encode()
	}

	return nil
}
