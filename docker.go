package main

import (
	"net/http"
	"path/filepath"

	"docker.io/go-docker"
	"docker.io/go-docker/api"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/urfave/cli"
)

// adapted from docker.NewEnvClient
func newDockerClient(c *cli.Context) (*docker.Client, error) {
	var client *http.Client
	if dockerCertPath := c.GlobalString(fDockerCertPath.Name); dockerCertPath != "" {
		options := tlsconfig.Options{
			CAFile:             filepath.Join(dockerCertPath, "ca.pem"),
			CertFile:           filepath.Join(dockerCertPath, "cert.pem"),
			KeyFile:            filepath.Join(dockerCertPath, "key.pem"),
			InsecureSkipVerify: c.GlobalString(fDockerTLSVerify.Name) == "",
		}
		tlsc, err := tlsconfig.Client(options)
		if err != nil {
			return nil, err
		}

		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsc,
			},
			CheckRedirect: docker.CheckRedirect,
		}
	}

	host := c.GlobalString(fDockerHost.Name)
	if host == "" {
		host = docker.DefaultDockerHost
	}
	version := c.GlobalString(fDockerAPIVersion.Name)
	if version == "" {
		version = api.DefaultVersion
	}

	dckr, err := docker.NewClient(host, version, client, nil)
	if err != nil {
		return dckr, err
	}

	return dckr, nil
}
