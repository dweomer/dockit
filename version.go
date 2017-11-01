package main

import (
	"fmt"

	docker "docker.io/go-docker/api"
	consul "github.com/hashicorp/consul/version"

	"github.com/urfave/cli"
)

const (
	// Version is the urfave/cli App.Version
	Version = "unknown"
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%v version %v (docker-api=%v, consul-api=%v)\n", c.App.HelpName, c.App.Version, docker.DefaultVersion, consul.Version)
	}

	app.Commands = append(app.Commands, cli.Command{
		Name:   "version",
		Usage:  "Print the version",
		Action: cli.ShowVersion,
	})
}
