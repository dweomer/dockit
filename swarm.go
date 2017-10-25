package main

import (
	"context"
	"errors"
	"fmt"

	docker "docker.io/go-docker"
	swarm "docker.io/go-docker/api/types/swarm"
	consul "github.com/hashicorp/consul/api"
	"github.com/urfave/cli"
)

func swarmInit(dcl *docker.Client, ccl *consul.Client, c *cli.Context) error {
	nodeID, err := dcl.SwarmInit(context.Background(), swarm.InitRequest{
		AdvertiseAddr:    c.String(fSwarmAdvertiseAddr.Name),
		AutoLockManagers: c.Bool(fSwarmAutolock.Name),
		Availability:     swarm.NodeAvailability(c.String(fSwarmAvailability.Name)),
		DataPathAddr:     c.String(fSwarmDataPathAddr.Name),
		ListenAddr:       c.String(fSwarmListenAddr.Name),
	})
	if err != nil {
		return err
	}

	swarm, err := dcl.SwarmInspect(context.Background())
	if err != nil {
		return err
	}

	node, _, err := dcl.NodeInspectWithRaw(context.Background(), nodeID)
	if err != nil {
		return err
	}

	if node.ManagerStatus != nil && node.ManagerStatus.Leader {
		_, err = ccl.KV().Put(&consul.KVPair{Key: "docker/swarm/id", Value: []byte(swarm.ClusterInfo.ID)}, nil)
		if err != nil {
			return err
		}

		_, err = ccl.KV().Put(&consul.KVPair{Key: "docker/swarm/leader", Value: []byte(node.ManagerStatus.Addr)}, nil)
		if err != nil {
			return err
		}

		_, err = ccl.KV().Put(&consul.KVPair{Key: "docker/swarm/join-token/manager", Value: []byte(swarm.JoinTokens.Manager)}, nil)
		if err != nil {
			return err
		}

		_, err = ccl.KV().Put(&consul.KVPair{Key: "docker/swarm/join-token/worker", Value: []byte(swarm.JoinTokens.Worker)}, nil)
		if err != nil {
			return err
		}
	} else {
		return errors.New("somehow not leader immediately after initializing the swarm")
	}

	return nil
}

func swarmJoin(dc *docker.Client, cc *consul.Client, c *cli.Context) error {
	tkey := fmt.Sprintf("docker/swarm/join-token/%v", c.Command.Name)
	tkvp, _, err := cc.KV().Get(tkey, nil)
	if err != nil {
		return err
	} else if tkvp == nil || len(tkvp.Value) == 0 {
		return fmt.Errorf("missing value for '%s'", tkey)
	}

	lkey := "docker/swarm/leader"
	lkvp, _, err := cc.KV().Get(lkey, nil)
	if err != nil {
		return err
	} else if lkvp == nil || len(lkvp.Value) == 0 {
		return fmt.Errorf("missing value for '%s'", lkey)
	}

	return dc.SwarmJoin(context.Background(), swarm.JoinRequest{
		AdvertiseAddr: c.String(fSwarmAdvertiseAddr.Name),
		Availability:  swarm.NodeAvailability(c.String(fSwarmAvailability.Name)),
		DataPathAddr:  c.String(fSwarmDataPathAddr.Name),
		ListenAddr:    c.String(fSwarmListenAddr.Name),
		JoinToken:     string(tkvp.Value),
		RemoteAddrs:   []string{string(lkvp.Value)},
	})
}
