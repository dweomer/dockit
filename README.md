# `sweomer` - Docker Swarm, with Consul

*Dweomer + Swarmer = Sweomer.*

Sweomer solves the Docker Swarm bootstrap problem by leveraging Consul to record and replicate the result of a race to initialize the Swarm. The first `manager` node to get the lock in Consul and populate the keys necessary for other managers and workers to join is the initial leader. After bootstrapping the Swarm is responsible for leadership and consensus via Swarm's built-in Raft implementation.

Sweomer is currently so simple that it will exit with a non-zero value at the first hint of a problem. It is assumed that Sweomer is executed in an environment where it is trivial to retry executions until success is achieved (i.e. `docker run --restart-policy=on-failure`).

Sweomer achives a rudimentary idempotency from invocation to invocation by checking Consul for the results of a prior successful execution and exiting cleanly if such already exists. The execution recipe runs something like this (with any error treated as a failure condition resulting in a non-zero exit code):
* Initialize the Consul client.
* Verify that the Consul agent is part of a cluster that has achieved quorum (aka check the Raft status) with a few configurable retries.
* Initialize the Docker client.
* Inspect the Docker engine.
  * If the node has any Swarm detail, exit cleanly - we're done.
* If this execution specifies a worker node, read the worker join token from Consul and attempt to join.
  * If the join is successful, exit cleanly - we're done.
* If this execution specifies a manager node, attempt to get a lock on the leadership key.
  * If we fail to get the lock, exit with a non-zero code.
  * While holding the lock: if there is no leadership detail in Consul, assume that this node is the leader and initialize the Swarm.
    * Write the Swarm details into Consul, release the lock, and exit cleanly - we're done.
  * While holding the lock: if there is leadership detail, release the lock, read the manager join token from Consul and attempt to join.
    * If the join is successful exit cleanly - we're done.

All successful executions result in the Swarm Node ID printed to standard out.

## Usage

### *`sweomer help`*

```
NAME:
   Sweomer - Docker Swarm, with Consul

USAGE:
   sweomer [global options] command [command options] [arguments...]

VERSION:
   unknown

COMMANDS:
     version  Print the version
     help, h  Shows a list of commands or help for one command
   Docker Swarm:
     manager  Swarm manager init/join, leveraging Consul
     worker   Swarm worker join, leveraging Consul

GLOBAL OPTIONS:
   --consul-cacert value           the Consul CA certificate file [$CONSUL_CACERT]
   --consul-capath value           the Consul CA certificate path [$CONSUL_CAPATH]
   --consul-client-cert value      the Consul client certificate [$CONSUL_CLIENT_CERT]
   --consul-client-key value       the Consul client key [$CONSUL_CLIENT_KEY]
   --consul-http-addr value        the Consul HTTP address [$CONSUL_HTTP_ADDR]
   --consul-http-auth value        the Consul HTTP authentication header [$CONSUL_HTTP_AUTH]
   --consul-http-ssl value         enable Consul over HTTPS? [$CONSUL_HTTP_SSL]
   --consul-http-ssl-verify value  enable Consul certificate verification? [$CONSUL_HTTP_SSL_VERIFY]
   --consul-http-token value       the Consul HTTP token [$CONSUL_HTTP_TOKEN]
   --consul-tls-server-name value  the Consul TLS/SNI server name [$CONSUL_TLS_SERVER_NAME]
   --docker-api-version value      (default: "1.33") [$DOCKER_API_VERSION]
   --docker-cert-path value        (default: /home/jacob/.docker) [$DOCKER_CERT_PATH]
   --docker-host value             as 'docker --host' [$DOCKER_HOST]
   --docker-tls-verify value       as 'docker --tlsverify' [$DOCKER_TLS_VERIFY]
   --help, -h                      show help
   --version, -v                   print the version
```

Sweomer leverages the official Docker and Consul client libraries. As such, the client-relevant environment variables leveraged by each runtime should be reusable as-is with command-line arguments having a runtime-specific prefix (or namespace), i.e. `sweomer --docker-host` instead of `docker --host` or `sweomer --consul-http-addr` instead of `consul --http-addr`. Additionally, the `sweomer manager` sub-command has `swarm`-prefixed arguments melding most of the arguments from `docker swarm init` and `docker swarm join` whereas the `sweomer worker` sub-command has the `swarm`-prefixed most of the arguments from `docker swarm join`. Both sub-commands also accept two new-to-Sweomer `--consul`-prefixed arguments for controlling how many times to retry checking for quorum and how long to wait between attempts.


### *`sweomer manager --help`*

```
NAME:
   sweomer manager - Swarm manager init/join, leveraging Consul

USAGE:
   sweomer manager [command options] [arguments...]

CATEGORY:
   Docker Swarm

OPTIONS:
   --consul-raft-retry-interval value  Time to wait between Consul Raft status reads (default: 5s)
   --consul-raft-retry-max value       Maximum number of Consul Raft status read attempts (default: 6)
   --swarm-advertise-addr value        Advertised address (format: <ip|interface>[:port])
   --swarm-autolock                    Enable manager autolocking (requiring an unlock key to start a stopped manager)
   --swarm-availability value          Availability of the node ("active"|"pause"|"drain") (default: "active")
   --swarm-data-path-addr value        Address or interface to use for data path traffic (format: <ip|interface>)
   --swarm-listen-addr value           Listen address (format: <ip|interface>[:port]) (default: "0.0.0.0:2377")
```

### *`sweomer worker --help`*

```
NAME:
   sweomer worker - Swarm worker join, leveraging Consul

USAGE:
   sweomer worker [command options] [arguments...]

CATEGORY:
   Docker Swarm

OPTIONS:
   --consul-raft-retry-interval value  Time to wait between Consul Raft status reads (default: 5s)
   --consul-raft-retry-max value       Maximum number of Consul Raft status read attempts (default: 6)
   --swarm-advertise-addr value        Advertised address (format: <ip|interface>[:port])
   --swarm-availability value          Availability of the node ("active"|"pause"|"drain") (default: "active")
   --swarm-data-path-addr value        Address or interface to use for data path traffic (format: <ip|interface>)
   --swarm-listen-addr value           Listen address (format: <ip|interface>[:port]) (default: "0.0.0.0:2377")
```

## Building

Sweomer is built with [go 1.9+](https://golang.org/dl/). Dependencies are managed with [Dep](https://github.com/golang/dep).

```
# install go 1.9 or better
# set a GOPATH, or don't
# install `dep`
go get -v github.com/golang/dep/cmd/dep
# get the source
go get -v -d github.com/dweomer/sweomer
cd $GOPATH/src/github.com/dweomer/sweomer
# get the dependencies into `./vendor`
dep ensure -v
# build and install into `$GOPATH/bin`
go install -v ./...
```

Always refer the `Dockerfile` for the most up-to-date recipe for building in a clean environment.

## Legalese
>The [MIT License](LICENSE) ([MIT](https://opensource.org/licenses/MIT))
>
> Copyright &copy; 2017 [Jacob Blain Christen](https://github.com/dweomer)
>
> Permission is hereby granted, free of charge, to any person obtaining a copy of
> this software and associated documentation files (the "Software"), to deal in
> the Software without restriction, including without limitation the rights to
> use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
> the Software, and to permit persons to whom the Software is furnished to do so,
> subject to the following conditions:
>
> The above copyright notice and this permission notice shall be included in all
> copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
> FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
> COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
> IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
> CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
