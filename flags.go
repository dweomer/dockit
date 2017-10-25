package main

import (
	"fmt"
	"os"
	"time"

	docker "docker.io/go-docker/api"
	consul "github.com/hashicorp/consul/api"
	"github.com/urfave/cli"
)

// CONSUL
var (
	fConsulCACert = cli.StringFlag{
		Name:   "consul-cacert",
		Usage:  "the Consul CA certificate file",
		EnvVar: consul.HTTPCAFile,
		Hidden: true,
	}
	fConsulCAPath = cli.StringFlag{
		Name:   "consul-capath",
		Usage:  "the Consul CA certificate path",
		EnvVar: consul.HTTPCAPath,
		Hidden: true,
	}
	fConsulClientCert = cli.StringFlag{
		Name:   "consul-client-cert",
		Usage:  "the Consul client certificate",
		EnvVar: consul.HTTPClientCert,
		Hidden: true,
	}
	fConsulClientKey = cli.StringFlag{
		Name:   "consul-client-key",
		Usage:  "the Consul client key",
		EnvVar: consul.HTTPClientKey,
		Hidden: true,
	}
	fConsulHTTPAddr = cli.StringFlag{
		Name:   "consul-http-addr",
		Usage:  "the Consul HTTP address",
		EnvVar: consul.HTTPAddrEnvName,
		Hidden: true,
	}
	fConsulHTTPAuth = cli.StringFlag{
		Name:   "consul-http-auth",
		Usage:  "the Consul HTTP authentication header",
		EnvVar: consul.HTTPAuthEnvName,
		Hidden: true,
	}
	fConsulHTTPToken = cli.StringFlag{
		Name:   "consul-http-token",
		Usage:  "the Consul HTTP token",
		EnvVar: consul.HTTPTokenEnvName,
		Hidden: true,
	}
	fConsulHTTPSSL = cli.StringFlag{
		Name:   "consul-http-ssl",
		Usage:  "enable Consul over HTTPS?",
		EnvVar: consul.HTTPSSLEnvName,
		Hidden: true,
	}
	fConsulHTTPSSLVerify = cli.StringFlag{
		Name:   "consul-http-ssl-verify",
		Usage:  "enable Consul certificate verification?",
		EnvVar: consul.HTTPSSLVerifyEnvName,
		Hidden: true,
	}
	fConsulTLSServerName = cli.StringFlag{
		Name:   "consul-tls-server-name",
		Usage:  "the Consul TLS/SNI server name",
		EnvVar: consul.HTTPTLSServerName,
		Hidden: true,
	}
	fConsulRaftRetryMax = cli.IntFlag{
		Name:  "consul-raft-retry-max",
		Usage: "Maximum number of Consul Raft status read attempts",
		Value: 6,
	}
	fConsulRaftRetryInterval = cli.DurationFlag{
		Name:  "consul-raft-retry-interval",
		Usage: "Time to wait between Consul Raft status reads",
		Value: 5 * time.Second,
	}
)

// DOCKER
var (
	fDockerAPIVersion = cli.StringFlag{
		Name:   "docker-api-version",
		EnvVar: "DOCKER_API_VERSION",
		Value:  docker.DefaultVersion,
	}
	fDockerCertPath = cli.StringFlag{
		Name:   "docker-cert-path",
		EnvVar: "DOCKER_CERT_PATH",
		Usage:  fmt.Sprintf("(default: %s/.docker)", os.Getenv("HOME")),
	}
	fDockerHost = cli.StringFlag{
		Name:   "docker-host",
		Usage:  "as 'docker --host'",
		EnvVar: "DOCKER_HOST",
	}
	fDockerTLSVerify = cli.StringFlag{
		Name:   "docker-tls-verify",
		Usage:  "as 'docker --tlsverify'",
		EnvVar: "DOCKER_TLS_VERIFY",
	}
)

// SWARM
var (
	fSwarmAdvertiseAddr = cli.StringFlag{
		Name:  "swarm-advertise-addr",
		Usage: "Advertised address (format: <ip|interface>[:port])",
	}
	fSwarmAutolock = cli.BoolFlag{
		Name:  "swarm-autolock",
		Usage: "Enable manager autolocking (requiring an unlock key to start a stopped manager)",
	}
	fSwarmAvailability = cli.StringFlag{
		Name:  "swarm-availability",
		Usage: "Availability of the node (\"active\"|\"pause\"|\"drain\")",
		Value: "active",
	}
	fSwarmCertExpiry = cli.DurationFlag{
		Name:   "swarm-cert-expiry",
		Usage:  "Validity period for node certificates (ns|us|ms|s|m|h)",
		Value:  90 * 24 * time.Hour,
		Hidden: true,
	}
	fSwarmDataPathAddr = cli.StringFlag{
		Name:  "swarm-data-path-addr",
		Usage: "Address or interface to use for data path traffic (format: <ip|interface>)",
	}
	fSwarmDispatcherHeartbeat = cli.DurationFlag{
		Name:   "swarm-dispatcher-heartbeat",
		Usage:  "Dispatcher heartbeat period (ns|us|ms|s|m|h)",
		Value:  5 * time.Second,
		Hidden: true,
	}
	fSwarmExternalCA = cli.StringSliceFlag{
		Name:   "swarm-external-ca",
		Usage:  "Specifications of one or more certificate signing endpoints",
		Hidden: true,
	}
	fSwarmForceNewCluster = cli.BoolFlag{
		Name:   "swarm-force-new-cluster",
		Usage:  "Force create a new cluster from current state",
		Hidden: true,
	}
	fSwarmListenAddr = cli.StringFlag{
		Name:  "swarm-listen-addr",
		Usage: "Listen address (format: <ip|interface>[:port])",
		Value: "0.0.0.0:2377",
	}
	fSwarmMaxSnapshots = cli.UintFlag{
		Name:   "swarm-max-snapshots",
		Usage:  "Number of additional Raft snapshots to retain",
		Hidden: true,
	}
	fSwarmSnapshotInterval = cli.UintFlag{
		Name:   "swarm-snapshot-interval",
		Usage:  "Number of log entries between Raft snapshots",
		Value:  10000,
		Hidden: true,
	}
	fSwarmTaskHistoryLimit = cli.Int64Flag{
		Name:   "swarm-task-history-limit",
		Usage:  "Task history retention limit",
		Value:  5,
		Hidden: true,
	}
)
