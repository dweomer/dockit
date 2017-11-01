package main

import (
	"fmt"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"github.com/urfave/cli"
)

// defaults/overrides logic adapted from Consul api.DefaultConfig()
func newConsulClient(c *cli.Context) (*consul.Client, error) {
	config := consul.DefaultConfig()

	if addr := c.GlobalString(fConsulHTTPAddr.Name); addr != "" {
		config.Address = addr
	}

	if token := c.GlobalString(fConsulHTTPToken.Name); token != "" {
		config.Token = token
	}

	if auth := c.GlobalString(fConsulHTTPAuth.Name); auth != "" {
		var username, password string
		if strings.Contains(auth, ":") {
			split := strings.SplitN(auth, ":", 2)
			username = split[0]
			password = split[1]
		} else {
			username = auth
		}

		config.HttpAuth = &consul.HttpBasicAuth{
			Username: username,
			Password: password,
		}
	}

	if ssl := c.GlobalString(fConsulHTTPSSL.Name); ssl != "" {
		enabled, err := strconv.ParseBool(ssl)
		if err != nil {
			fmt.Fprintf(c.App.ErrWriter, "Failed to parse %s: '%s', skipped.", fConsulHTTPSSL.Name, err)
		}
		if enabled {
			config.Scheme = "https"
		}
	}

	if v := c.GlobalString(fConsulTLSServerName.Name); v != "" {
		config.TLSConfig.Address = v
	}
	if v := c.GlobalString(fConsulCACert.Name); v != "" {
		config.TLSConfig.CAFile = v
	}
	if v := c.GlobalString(fConsulCAPath.Name); v != "" {
		config.TLSConfig.CAPath = v
	}
	if v := c.GlobalString(fConsulClientCert.Name); v != "" {
		config.TLSConfig.CertFile = v
	}
	if v := c.GlobalString(fConsulClientKey.Name); v != "" {
		config.TLSConfig.KeyFile = v
	}
	if v := c.GlobalString(fConsulHTTPSSLVerify.Name); v != "" {
		verify, err := strconv.ParseBool(v)
		if err != nil {
			fmt.Fprintf(c.App.ErrWriter, "Failed to parse %s: '%s', skipped.", fConsulHTTPSSLVerify.Name, err)
		}
		if !verify {
			config.TLSConfig.InsecureSkipVerify = true
		}
	}

	return consul.NewClient(config)
}
